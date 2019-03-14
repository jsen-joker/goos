/*
 * Copyright (C) 2018 the original author or authors.
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

package org.jsen.goos.locator;

import org.jsen.goos.client.ConfigService;
import org.jsen.goos.client.exception.GoosClientException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.config.YamlPropertiesFactoryBean;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.util.StringUtils;

import java.io.StringReader;
import java.util.*;

/**
 * @author xiaojing
 */
public class GoosPropertySourceBuilder {

	private static final Logger logger = LoggerFactory
			.getLogger(GoosPropertySourceBuilder.class);

	private ConfigService configService;
	private long timeout;

	public GoosPropertySourceBuilder() {
	}

	public GoosPropertySourceBuilder(ConfigService configService, long timeout) {
		this.configService = configService;
		this.timeout = timeout;
	}

	public long getTimeout() {
		return timeout;
	}

	public void setTimeout(long timeout) {
		this.timeout = timeout;
	}

	public ConfigService getConfigService() {
		return configService;
	}

	public void setConfigService(ConfigService configService) {
		this.configService = configService;
	}

	/**
	 * @param dataId Nacos dataId
	 * @param group Nacos group
	 */
	GoosPropertySource build(String dataId, String group, String fileExtension) {
		Properties p = loadNacosData(dataId, group, fileExtension);
		if (p == null) {
			return null;
		}
		return new GoosPropertySource(dataId, propertiesToMap(p), new Date());
	}

	private Properties loadNacosData(String dataId, String group, String fileExtension) {
		String data = null;
		try {
			data = configService.getConfig(dataId, group, timeout);
			if (!StringUtils.isEmpty(data)) {
				logger.info(String.format("Loading nacos data, dataId: '%s', group: '%s'",
						dataId, group));

				if (fileExtension.equalsIgnoreCase("properties")) {
					Properties properties = new Properties();

					properties.load(new StringReader(data));
					return properties;
				}
				else if (fileExtension.equalsIgnoreCase("yaml")
						|| fileExtension.equalsIgnoreCase("yml")) {
					YamlPropertiesFactoryBean yamlFactory = new YamlPropertiesFactoryBean();
					yamlFactory.setResources(new ByteArrayResource(data.getBytes()));
					return yamlFactory.getObject();
				}

			}
		}
		catch (GoosClientException e) {
			logger.error("get data from Nacos error,dataId:{}, ", dataId, e);
		}
		catch (Exception e) {
			logger.error("parse data from Nacos error,dataId:{},data:{},", dataId, data,
					e);
		}
		return null;
	}

	@SuppressWarnings("unchecked")
	private Map<String, Object> propertiesToMap(Properties properties) {
		Map<String, Object> result = new HashMap<>(16);
		Enumeration<String> keys = (Enumeration<String>) properties.propertyNames();
		while (keys.hasMoreElements()) {
			String key = keys.nextElement();
			Object value = properties.getProperty(key);
			if (value != null) {
				result.put(key, ((String) value).trim());
			}
			else {
				result.put(key, null);
			}
		}
		return result;
	}
}
