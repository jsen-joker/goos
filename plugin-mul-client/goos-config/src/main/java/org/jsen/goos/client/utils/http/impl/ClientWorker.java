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
package org.jsen.goos.client.utils.http.impl;

import org.jsen.goos.client.common.GroupKey;
import org.jsen.goos.client.exception.GoosClientException;
import org.jsen.goos.client.filter.impl.ConfigFilterChainManager;
import org.jsen.goos.client.listener.Listener;
import org.jsen.goos.client.monitor.MetricsMonitor;
import org.jsen.goos.client.utils.*;
import org.jsen.goos.client.utils.http.Constants;
import org.jsen.goos.client.utils.http.HttpAgent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.URLDecoder;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicReference;

import static org.jsen.goos.client.utils.http.Constants.LINE_SEPARATOR;
import static org.jsen.goos.client.utils.http.Constants.WORD_SEPARATOR;


/**
 * wait for fit goos api
 *
 * @author Goos
 */
public class ClientWorker {

    private static final Logger LOGGER = LoggerFactory.getLogger(ClientWorker.class);

    public void addListeners(String dataId, String group, List<? extends Listener> listeners) {
        group = null2defaultGroup(group);
        CacheData cache = addCacheDataIfAbsent(dataId, group);
        for (Listener listener : listeners) {
            cache.addListener(listener);
        }
    }

    public void removeListener(String dataId, String group, Listener listener) {
        group = null2defaultGroup(group);
        CacheData cache = getCache(dataId, group);
        if (null != cache) {
            cache.removeListener(listener);
            if (cache.getListeners().isEmpty()) {
                removeCache(dataId, group);
            }
        }
    }

    public void addTenantListeners(String dataId, String group, List<? extends Listener> listeners) {
        group = null2defaultGroup(group);
        String tenant = agent.getTenant();
        CacheData cache = addCacheDataIfAbsent(dataId, group, tenant);
        for (Listener listener : listeners) {
            cache.addListener(listener);
        }
    }

    public void removeTenantListener(String dataId, String group, Listener listener) {
        group = null2defaultGroup(group);
        String tenant = agent.getTenant();
        CacheData cache = getCache(dataId, group, tenant);
        if (null != cache) {
            cache.removeListener(listener);
            if (cache.getListeners().isEmpty()) {
                removeCache(dataId, group, tenant);
            }
        }
    }

//    @SuppressFBWarnings("JLM_JSR166_UTILCONCURRENT_MONITORENTER")
    void removeCache(String dataId, String group) {
        String groupKey = GroupKey.getKey(dataId, group);
        synchronized (cacheMap) {
            Map<String, CacheData> copy = new HashMap<String, CacheData>(cacheMap.get());
            copy.remove(groupKey);
            cacheMap.set(copy);
        }
        LOGGER.info("[{}] [unsubscribe] {}", agent.getName(), groupKey);

        MetricsMonitor.getListenConfigCountMonitor().set(cacheMap.get().size());
    }

//    @SuppressFBWarnings("JLM_JSR166_UTILCONCURRENT_MONITORENTER")
    void removeCache(String dataId, String group, String tenant) {
        String groupKey = GroupKey.getKeyTenant(dataId, group, tenant);
        synchronized (cacheMap) {
            Map<String, CacheData> copy = new HashMap<String, CacheData>(cacheMap.get());
            copy.remove(groupKey);
            cacheMap.set(copy);
        }
        LOGGER.info("[{}] [unsubscribe] {}", agent.getName(), groupKey);

        MetricsMonitor.getListenConfigCountMonitor().set(cacheMap.get().size());
    }

//    @SuppressFBWarnings("JLM_JSR166_UTILCONCURRENT_MONITORENTER")
    public CacheData addCacheDataIfAbsent(String dataId, String group) {
        CacheData cache = getCache(dataId, group);
        if (null != cache) {
            return cache;
        }

        String key = GroupKey.getKey(dataId, group);
        cache = new CacheData(configFilterChainManager, agent.getName(), dataId, group);

        synchronized (cacheMap) {
            CacheData cacheFromMap = getCache(dataId, group);
            // multiple listeners on the same dataid+group and race condition,so double check again
            //other listener thread beat me to set to cacheMap
            if (null != cacheFromMap) {
                cache = cacheFromMap;
                //reset so that server not hang this check
                cache.setInitializing(true);
            } else {
                int taskId = cacheMap.get().size() / (int)ParamUtil.getPerTaskConfigSize();
                cache.setTaskId(taskId);
            }

            Map<String, CacheData> copy = new HashMap<String, CacheData>(cacheMap.get());
            copy.put(key, cache);
            cacheMap.set(copy);
        }

        LOGGER.info("[{}] [subscribe] {}", agent.getName(), key);

        MetricsMonitor.getListenConfigCountMonitor().set(cacheMap.get().size());

        return cache;
    }

//    @SuppressFBWarnings("JLM_JSR166_UTILCONCURRENT_MONITORENTER")
    public CacheData addCacheDataIfAbsent(String dataId, String group, String tenant) {
        CacheData cache = getCache(dataId, group, tenant);
        if (null != cache) {
            return cache;
        }
        String key = GroupKey.getKeyTenant(dataId, group, tenant);
        cache = new CacheData(configFilterChainManager, agent.getName(), dataId, group, tenant);
        synchronized (cacheMap) {
            CacheData cacheFromMap = getCache(dataId, group, tenant);
            // multiple listeners on the same dataid+group and race condition,so
            // double check again
            // other listener thread beat me to set to cacheMap
            if (null != cacheFromMap) {
                cache = cacheFromMap;
                // reset so that server not hang this check
                cache.setInitializing(true);
            }

            Map<String, CacheData> copy = new HashMap<String, CacheData>(cacheMap.get());
            copy.put(key, cache);
            cacheMap.set(copy);
        }
        LOGGER.info("[{}] [subscribe] {}", agent.getName(), key);

        MetricsMonitor.getListenConfigCountMonitor().set(cacheMap.get().size());

        return cache;
    }

    public CacheData getCache(String dataId, String group) {
        return getCache(dataId, group, TenantUtil.getUserTenant());
    }

    public CacheData getCache(String dataId, String group, String tenant) {
        if (null == dataId || null == group) {
            throw new IllegalArgumentException();
        }
        return cacheMap.get().get(GroupKey.getKeyTenant(dataId, group, tenant));
    }

    public String getServerConfig(String dataId, String group, String tenant, long readTimeout)
        throws GoosClientException {
        if (StringUtils.isBlank(group)) {
            group = Constants.DEFAULT_GROUP;
        }

        HttpSimpleClient.HttpResult result = null;
        try {
            List<String> params = null;
            if (StringUtils.isBlank(tenant)) {
                params = Arrays.asList("dataId", dataId, "group", group);
            } else {
                params = Arrays.asList("dataId", dataId, "group", group, "tenant", tenant);
            }
            result = agent.httpGet(Constants.CONFIG_CONTROLLER_PATH, null, params, agent.getEncode(), readTimeout);
        } catch (IOException e) {
            String message = String.format(
                "[%s] [sub-server] get server config exception, dataId=%s, group=%s, tenant=%s", agent.getName(),
                dataId, group, tenant);
            LOGGER.error(message, e);
            throw new GoosClientException(GoosClientException.SERVER_ERROR, e.getMessage());
        }

        switch (result.code) {
            case HttpURLConnection.HTTP_OK:
                LocalConfigInfoProcessor.saveSnapshot(agent.getName(), dataId, group, tenant, result.content);
                return result.content;
            case HttpURLConnection.HTTP_NOT_FOUND:
                LocalConfigInfoProcessor.saveSnapshot(agent.getName(), dataId, group, tenant, null);
                return null;
            case HttpURLConnection.HTTP_CONFLICT: {
                LOGGER.error(
                    "[{}] [sub-server-error] get server config being modified concurrently, dataId={}, group={}, "
                        + "tenant={}", agent.getName(), dataId, group, tenant);
                throw new GoosClientException(GoosClientException.CONFLICT,
                    "data being modified, dataId=" + dataId + ",group=" + group + ",tenant=" + tenant);
            }
            case HttpURLConnection.HTTP_FORBIDDEN: {
                LOGGER.error("[{}] [sub-server-error] no right, dataId={}, group={}, tenant={}", agent.getName(), dataId,
                    group, tenant);
                throw new GoosClientException(result.code, result.content);
            }
            default: {
                LOGGER.error("[{}] [sub-server-error]  dataId={}, group={}, tenant={}, code={}", agent.getName(), dataId,
                    group, tenant, result.code);
                throw new GoosClientException(result.code,
                    "http error, code=" + result.code + ",dataId=" + dataId + ",group=" + group + ",tenant=" + tenant);
            }
        }
    }

    private void checkLocalConfig(CacheData cacheData) {
        final String dataId = cacheData.dataId;
        final String group = cacheData.group;
        final String tenant = cacheData.tenant;
        File path = LocalConfigInfoProcessor.getFailoverFile(agent.getName(), dataId, group, tenant);

        // 没有 -> 有
        if (!cacheData.isUseLocalConfigInfo() && path.exists()) {
            String content = LocalConfigInfoProcessor.getFailover(agent.getName(), dataId, group, tenant);
            String md5 = MD5.getInstance().getMD5String(content);
            cacheData.setUseLocalConfigInfo(true);
            cacheData.setLocalConfigInfoVersion(path.lastModified());
            cacheData.setContent(content);

            LOGGER.warn("[{}] [failover-change] failover file created. dataId={}, group={}, tenant={}, md5={}, content={}",
                agent.getName(), dataId, group, tenant, md5, ContentUtils.truncateContent(content));
            return;
        }

        // 有 -> 没有。不通知业务监听器，从server拿到配置后通知。
        if (cacheData.isUseLocalConfigInfo() && !path.exists()) {
            cacheData.setUseLocalConfigInfo(false);
            LOGGER.warn("[{}] [failover-change] failover file deleted. dataId={}, group={}, tenant={}", agent.getName(),
                dataId, group, tenant);
            return;
        }

        // 有变更
        if (cacheData.isUseLocalConfigInfo() && path.exists()
            && cacheData.getLocalConfigInfoVersion() != path.lastModified()) {
            String content = LocalConfigInfoProcessor.getFailover(agent.getName(), dataId, group, tenant);
            String md5 = MD5.getInstance().getMD5String(content);
            cacheData.setUseLocalConfigInfo(true);
            cacheData.setLocalConfigInfoVersion(path.lastModified());
            cacheData.setContent(content);
            LOGGER.warn("[{}] [failover-change] failover file changed. dataId={}, group={}, tenant={}, md5={}, content={}",
                agent.getName(), dataId, group, tenant, md5, ContentUtils.truncateContent(content));
            return;
        }
    }

    private String null2defaultGroup(String group) {
        return (null == group) ? Constants.DEFAULT_GROUP : group.trim();
    }

    public void checkConfigInfo() {
        // 分任务
        int listenerSize = cacheMap.get().size();
        // 向上取整为批数
        int longingTaskCount = (int) Math.ceil(listenerSize / ParamUtil.getPerTaskConfigSize());
        if (longingTaskCount > currentLongingTaskCount) {
            for (int i = (int)currentLongingTaskCount; i < longingTaskCount; i++) {
                // 要判断任务是否在执行 这块需要好好想想。 任务列表现在是无序的。变化过程可能有问题
                executorService.execute(new LongPollingRunnable(i));
            }
            currentLongingTaskCount = longingTaskCount;
        }
    }

    /**
     * 从Server获取值变化了的DataID列表。返回的对象里只有dataId和group是有效的。 保证不返回NULL。
     */
    List<String> checkUpdateDataIds(List<CacheData> cacheDatas, List<String> inInitializingCacheList) {
        StringBuilder sb = new StringBuilder();
        for (CacheData cacheData : cacheDatas) {
            if (!cacheData.isUseLocalConfigInfo()) {
                sb.append(cacheData.dataId).append(WORD_SEPARATOR);
                sb.append(cacheData.group).append(WORD_SEPARATOR);
                if (StringUtils.isBlank(cacheData.tenant)) {
                    sb.append(cacheData.getMd5()).append(LINE_SEPARATOR);
                } else {
                    sb.append(cacheData.getMd5()).append(WORD_SEPARATOR);
                    sb.append(cacheData.getTenant()).append(LINE_SEPARATOR);
                }
                if (cacheData.isInitializing()) {
                    // cacheData 首次出现在cacheMap中&首次check更新
                    inInitializingCacheList
                        .add(GroupKey.getKeyTenant(cacheData.dataId, cacheData.group, cacheData.tenant));
                }
            }
        }
        boolean isInitializingCacheList = !inInitializingCacheList.isEmpty();
        return checkUpdateConfigStr(sb.toString(), isInitializingCacheList);
    }

    /**
     * 从Server获取值变化了的DataID列表。返回的对象里只有dataId和group是有效的。 保证不返回NULL。
     */
    List<String> checkUpdateConfigStr(String probeUpdateString, boolean isInitializingCacheList) {

        List<String> params = Arrays.asList(Constants.PROBE_MODIFY_REQUEST, probeUpdateString);
        long timeout = TimeUnit.SECONDS.toMillis(30L);

        List<String> headers = new ArrayList<String>(2);
        headers.add("Long-Pulling-Timeout");
        headers.add("" + timeout);

        // told server do not hang me up if new initializing cacheData added in
        if (isInitializingCacheList) {
            headers.add("Long-Pulling-Timeout-No-Hangup");
            headers.add("true");
        }

        if (StringUtils.isBlank(probeUpdateString)) {
            return Collections.emptyList();
        }

        try {
            HttpSimpleClient.HttpResult result = agent.httpPost(Constants.CONFIG_CONTROLLER_PATH + "/listener", headers, params,
                agent.getEncode(), timeout);

            if (HttpURLConnection.HTTP_OK == result.code) {
                setHealthServer(true);
                return parseUpdateDataIdResponse(result.content);
            } else {
                setHealthServer(false);
                LOGGER.error("[{}] [check-update] get changed dataId error, code: {}", agent.getName(), result.code);
            }
        } catch (IOException e) {
            setHealthServer(false);
            LOGGER.error("[" + agent.getName() + "] [check-update] get changed dataId exception", e);
        }
        return Collections.emptyList();
    }

    /**
     * 从HTTP响应拿到变化的groupKey。保证不返回NULL。
     */
    private List<String> parseUpdateDataIdResponse(String response) {
        if (StringUtils.isBlank(response)) {
            return Collections.emptyList();
        }

        try {
            response = URLDecoder.decode(response, "UTF-8");
        } catch (Exception e) {
            LOGGER.error("[" + agent.getName() + "] [polling-resp] decode modifiedDataIdsString error", e);
        }

        List<String> updateList = new LinkedList<String>();

        for (String dataIdAndGroup : response.split(LINE_SEPARATOR)) {
            if (!StringUtils.isBlank(dataIdAndGroup)) {
                String[] keyArr = dataIdAndGroup.split(WORD_SEPARATOR);
                String dataId = keyArr[0];
                String group = keyArr[1];
                if (keyArr.length == 2) {
                    updateList.add(GroupKey.getKey(dataId, group));
                    LOGGER.info("[{}] [polling-resp] config changed. dataId={}, group={}", agent.getName(), dataId, group);
                } else if (keyArr.length == 3) {
                    String tenant = keyArr[2];
                    updateList.add(GroupKey.getKeyTenant(dataId, group, tenant));
                    LOGGER.info("[{}] [polling-resp] config changed. dataId={}, group={}, tenant={}", agent.getName(),
                        dataId, group, tenant);
                } else {
                    LOGGER.error("[{}] [polling-resp] invalid dataIdAndGroup error {}", agent.getName(), dataIdAndGroup);
                }
            }
        }
        return updateList;
    }

    @SuppressWarnings("PMD.ThreadPoolCreationRule")
    public ClientWorker(final HttpAgent agent, final ConfigFilterChainManager configFilterChainManager) {
        this.agent = agent;
        this.configFilterChainManager = configFilterChainManager;

        executor = Executors.newScheduledThreadPool(1, new ThreadFactory() {
            @Override
            public Thread newThread(Runnable r) {
                Thread t = new Thread(r);
                t.setName("com.alibaba.nacos.client.Worker." + agent.getName());
                t.setDaemon(true);
                return t;
            }
        });

        executorService = Executors.newCachedThreadPool(new ThreadFactory() {
            @Override
            public Thread newThread(Runnable r) {
                Thread t = new Thread(r);
                t.setName("com.alibaba.nacos.client.Worker.longPolling" + agent.getName());
                t.setDaemon(true);
                return t;
            }
        });

        executor.scheduleWithFixedDelay(new Runnable() {
            @Override
            public void run() {
                try {
                    checkConfigInfo();
                } catch (Throwable e) {
                    LOGGER.error("[" + agent.getName() + "] [sub-check] rotate check error", e);
                }
            }
        }, 1L, 10L, TimeUnit.MILLISECONDS);
    }

    class LongPollingRunnable implements Runnable {
        private int taskId;

        public LongPollingRunnable(int taskId) {
            this.taskId = taskId;
        }

        @Override
        public void run() {
            try {
                List<CacheData> cacheDatas = new ArrayList<CacheData>();
                // check failover config
                for (CacheData cacheData : cacheMap.get().values()) {
                    if (cacheData.getTaskId() == taskId) {
                        cacheDatas.add(cacheData);
                        try {
                            checkLocalConfig(cacheData);
                            if (cacheData.isUseLocalConfigInfo()) {
                                cacheData.checkListenerMd5();
                            }
                        } catch (Exception e) {
                            LOGGER.error("get local config info error", e);
                        }
                    }
                }

                List<String> inInitializingCacheList = new ArrayList<String>();
                // check server config
                List<String> changedGroupKeys = checkUpdateDataIds(cacheDatas, inInitializingCacheList);

                for (String groupKey : changedGroupKeys) {
                    String[] key = GroupKey.parseKey(groupKey);
                    String dataId = key[0];
                    String group = key[1];
                    String tenant = null;
                    if (key.length == 3) {
                        tenant = key[2];
                    }
                    try {
                        String content = getServerConfig(dataId, group, tenant, 3000L);
                        CacheData cache = cacheMap.get().get(GroupKey.getKeyTenant(dataId, group, tenant));
                        cache.setContent(content);
                        LOGGER.info("[{}] [data-received] dataId={}, group={}, tenant={}, md5={}, content={}",
                            agent.getName(), dataId, group, tenant, cache.getMd5(),
                            ContentUtils.truncateContent(content));
                    } catch (GoosClientException ioe) {
                        String message = String.format(
                            "[%s] [get-update] get changed config exception. dataId=%s, group=%s, tenant=%s",
                            agent.getName(), dataId, group, tenant);
                        LOGGER.error(message, ioe);
                    }
                }
                for (CacheData cacheData : cacheDatas) {
                    if (!cacheData.isInitializing() || inInitializingCacheList
                        .contains(GroupKey.getKeyTenant(cacheData.dataId, cacheData.group, cacheData.tenant))) {
                        cacheData.checkListenerMd5();
                        cacheData.setInitializing(false);
                    }
                }
                inInitializingCacheList.clear();
            } catch (Throwable e) {
                LOGGER.error("longPolling error", e);
            } finally {
                executorService.execute(this);
            }
        }
    }

    // =================

    public boolean isHealthServer() {
        return isHealthServer;
    }

    private void setHealthServer(boolean isHealthServer) {
        this.isHealthServer = isHealthServer;
    }

    final ScheduledExecutorService executor;
    final ExecutorService executorService;
    /**
     * groupKey -> cacheData
     */
    AtomicReference<Map<String, CacheData>> cacheMap = new AtomicReference<Map<String, CacheData>>(
        new HashMap<String, CacheData>());

    HttpAgent agent;
    ConfigFilterChainManager configFilterChainManager;
    private boolean isHealthServer = true;
    private double currentLongingTaskCount = 0;
}
