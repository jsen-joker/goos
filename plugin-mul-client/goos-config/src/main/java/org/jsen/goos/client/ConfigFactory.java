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
import org.jsen.goos.client.utils.http.impl.PropertyKeyConst;

import java.lang.reflect.Constructor;
import java.util.Properties;

/**
 * Config Factory
 *
 * @author Nacos
 */
public class ConfigFactory {

    /**
     * Create Config
     *
     * @param properties init param
     * @return ConfigService
     * @throws GoosConfigService Exception
     */
    public static ConfigService createConfigService(Properties properties) throws GoosClientException {
        try {
            Class<?> driverImplClass = Class.forName("org.jsen.goos.client.GoosConfigService");
            Constructor constructor = driverImplClass.getConstructor(Properties.class);
            ConfigService vendorImpl = (ConfigService) constructor.newInstance(properties);
            return vendorImpl;
        } catch (Throwable e) {
            throw new GoosClientException(-400, e.getMessage());
        }
    }

    /**
     * Create Config
     *
     * @param serverAddr serverList
     * @return Config
     * @throws ConfigService Exception
     */
    public static ConfigService createConfigService(String serverAddr) throws GoosClientException {
        Properties properties = new Properties();
        properties.put(PropertyKeyConst.SERVER_ADDR, serverAddr);
        return createConfigService(properties);
    }

}
