/*
 * Copyright 1999-2018 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.jsen.goos.client;

import org.jsen.goos.client.exception.GoosClientException;
import org.jsen.goos.client.filter.impl.ConfigFilterChainManager;
import org.jsen.goos.client.filter.impl.ConfigRequest;
import org.jsen.goos.client.filter.impl.ConfigResponse;
import org.jsen.goos.client.listener.Listener;
import org.jsen.goos.client.utils.ContentUtils;
import org.jsen.goos.client.utils.ParamUtils;
import org.jsen.goos.client.utils.StringUtils;
import org.jsen.goos.client.utils.TenantUtil;
import org.jsen.goos.client.utils.http.Constants;
import org.jsen.goos.client.utils.http.HttpAgent;
import org.jsen.goos.client.utils.http.MetricsHttpAgent;
import org.jsen.goos.client.utils.http.ServerHttpAgent;
import org.jsen.goos.client.utils.http.impl.ClientWorker;
import org.jsen.goos.client.utils.http.impl.HttpSimpleClient;
import org.jsen.goos.client.utils.http.impl.LocalConfigInfoProcessor;
import org.jsen.goos.client.utils.http.impl.PropertyKeyConst;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.HttpURLConnection;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Properties;

/**
 * Config Impl
 *
 * @author Nacos
 */
@SuppressWarnings("PMD.ServiceOrDaoClassShouldEndWithImplRule")
public class GoosConfigService implements ConfigService {

    private static final Logger LOGGER = LoggerFactory.getLogger(GoosConfigService.class);

    private final long POST_TIMEOUT = 3000L;
    /**
     * http agent
     */
    private HttpAgent agent;
    /**
     * longpolling
     */
    private ClientWorker worker;
    private String namespace;
    private String encode;
    private ConfigFilterChainManager configFilterChainManager = new ConfigFilterChainManager();

    public GoosConfigService(Properties properties) throws GoosClientException {
        String encodeTmp = properties.getProperty(PropertyKeyConst.ENCODE);
        if (StringUtils.isBlank(encodeTmp)) {
            encode = Constants.ENCODE;
        } else {
            encode = encodeTmp.trim();
        }
        String namespaceTmp = properties.getProperty(PropertyKeyConst.NAMESPACE);
        if (StringUtils.isBlank(namespaceTmp)) {
            namespace = TenantUtil.getUserTenant();
            properties.put(PropertyKeyConst.NAMESPACE, namespace);
        } else {
            namespace = namespaceTmp;
            properties.put(PropertyKeyConst.NAMESPACE, namespace);
        }
        agent = new MetricsHttpAgent(new ServerHttpAgent(properties));
        agent.start();
        worker = new ClientWorker(agent, configFilterChainManager);
    }

    @Override
    public String getConfig(String dataId, String group, long timeoutMs) throws GoosClientException {
        return getConfigInner(namespace, dataId, group, timeoutMs);
    }

    @Override
    public void addListener(String dataId, String group, Listener listener) throws GoosClientException {
        worker.addTenantListeners(dataId, group, Arrays.asList(listener));
    }

    @Override
    public boolean publishConfig(String dataId, String group, String content) throws GoosClientException {
        return publishConfigInner(namespace, dataId, group, null, null, null, content);
    }

    @Override
    public boolean removeConfig(String dataId, String group) throws GoosClientException {
        return removeConfigInner(namespace, dataId, group, null);
    }

    @Override
    public void removeListener(String dataId, String group, Listener listener) {
        worker.removeTenantListener(dataId, group, listener);
    }

    private String getConfigInner(String tenant, String dataId, String group, long timeoutMs) throws GoosClientException {
        group = null2defaultGroup(group);
        ParamUtils.checkKeyParam(dataId, group);
        ConfigResponse cr = new ConfigResponse();

        cr.setDataId(dataId);
        cr.setTenant(tenant);
        cr.setGroup(group);

        // 优先使用本地配置
        String content = LocalConfigInfoProcessor.getFailover(agent.getName(), dataId, group, tenant);
        if (content != null) {
            LOGGER.warn("[{}] [get-config] get failover ok, dataId={}, group={}, tenant={}, config={}", agent.getName(),
                dataId, group, tenant, ContentUtils.truncateContent(content));
            cr.setContent(content);
            configFilterChainManager.doFilter(null, cr);
            content = cr.getContent();
            return content;
        }

        try {
            content = worker.getServerConfig(dataId, group, tenant, timeoutMs);

            cr.setContent(content);
            configFilterChainManager.doFilter(null, cr);
            content = cr.getContent();

            return content;
        } catch (GoosClientException ioe) {
            if (GoosClientException.NO_RIGHT == ioe.getErrCode()) {
                throw ioe;
            }
            LOGGER.warn("[{}] [get-config] get from server error, dataId={}, group={}, tenant={}, msg={}",
                agent.getName(), dataId, group, tenant, ioe.toString());
        }

        LOGGER.warn("[{}] [get-config] get snapshot ok, dataId={}, group={}, tenant={}, config={}", agent.getName(),
            dataId, group, tenant, ContentUtils.truncateContent(content));
        content = LocalConfigInfoProcessor.getSnapshot(agent.getName(), dataId, group, tenant);
        cr.setContent(content);
        configFilterChainManager.doFilter(null, cr);
        content = cr.getContent();
        return content;
    }

    private String null2defaultGroup(String group) {
        return (null == group) ? Constants.DEFAULT_GROUP : group.trim();
    }

    private boolean removeConfigInner(String tenant, String dataId, String group, String tag) throws GoosClientException {
        group = null2defaultGroup(group);
        ParamUtils.checkKeyParam(dataId, group);
        String url = Constants.CONFIG_CONTROLLER_PATH;
        List<String> params = new ArrayList<String>();
        params.add("dataId");
        params.add(dataId);
        params.add("group");
        params.add(group);
        if (StringUtils.isNotEmpty(tenant)) {
            params.add("tenant");
            params.add(tenant);
        }
        if (StringUtils.isNotEmpty(tag)) {
            params.add("tag");
            params.add(tag);
        }
        HttpSimpleClient.HttpResult result = null;
        try {
            result = agent.httpDelete(url, null, params, encode, POST_TIMEOUT);
        } catch (IOException ioe) {
            LOGGER.warn("[remove] error, " + dataId + ", " + group + ", " + tenant + ", msg: " + ioe.toString());
            return false;
        }

        if (HttpURLConnection.HTTP_OK == result.code) {
            LOGGER.info("[{}] [remove] ok, dataId={}, group={}, tenant={}", agent.getName(), dataId, group, tenant);
            return true;
        } else if (HttpURLConnection.HTTP_FORBIDDEN == result.code) {
            LOGGER.warn("[{}] [remove] error, dataId={}, group={}, tenant={}, code={}, msg={}", agent.getName(), dataId,
                group, tenant, result.code, result.content);
            throw new GoosClientException(result.code, result.content);
        } else {
            LOGGER.warn("[{}] [remove] error, dataId={}, group={}, tenant={}, code={}, msg={}", agent.getName(), dataId,
                group, tenant, result.code, result.content);
            return false;
        }
    }

    private boolean publishConfigInner(String tenant, String dataId, String group, String tag, String appName,
                                       String betaIps, String content) throws GoosClientException {
        group = null2defaultGroup(group);
        ParamUtils.checkParam(dataId, group, content);

        ConfigRequest cr = new ConfigRequest();
        cr.setDataId(dataId);
        cr.setTenant(tenant);
        cr.setGroup(group);
        cr.setContent(content);
        configFilterChainManager.doFilter(cr, null);
        content = cr.getContent();

        String url = Constants.CONFIG_CONTROLLER_PATH;
        List<String> params = new ArrayList<String>();
        params.add("dataId");
        params.add(dataId);
        params.add("group");
        params.add(group);
        params.add("content");
        params.add(content);
        if (StringUtils.isNotEmpty(tenant)) {
            params.add("tenant");
            params.add(tenant);
        }
        if (StringUtils.isNotEmpty(appName)) {
            params.add("appName");
            params.add(appName);
        }
        if (StringUtils.isNotEmpty(tag)) {
            params.add("tag");
            params.add(tag);
        }

        List<String> headers = new ArrayList<String>();
        if (StringUtils.isNotEmpty(betaIps)) {
            headers.add("betaIps");
            headers.add(betaIps);
        }

        HttpSimpleClient.HttpResult result = null;
        try {
            result = agent.httpPost(url, headers, params, encode, POST_TIMEOUT);
        } catch (IOException ioe) {
            LOGGER.warn("[{}] [publish-single] exception, dataId={}, group={}, msg={}", agent.getName(), dataId,
                group, ioe.toString());
            return false;
        }

        if (HttpURLConnection.HTTP_OK == result.code) {
            LOGGER.info("[{}] [publish-single] ok, dataId={}, group={}, tenant={}, config={}", agent.getName(), dataId,
                group, tenant, ContentUtils.truncateContent(content));
            return true;
        } else if (HttpURLConnection.HTTP_FORBIDDEN == result.code) {
            LOGGER.warn("[{}] [publish-single] error, dataId={}, group={}, tenant={}, code={}, msg={}", agent.getName(),
                dataId, group, tenant, result.code, result.content);
            throw new GoosClientException(result.code, result.content);
        } else {
            LOGGER.warn("[{}] [publish-single] error, dataId={}, group={}, tenant={}, code={}, msg={}", agent.getName(),
                dataId, group, tenant, result.code, result.content);
            return false;
        }

    }

    @Override
    public String getServerStatus() {
        if (worker.isHealthServer()) {
            return "UP";
        } else {
            return "DOWN";
        }
    }

}
