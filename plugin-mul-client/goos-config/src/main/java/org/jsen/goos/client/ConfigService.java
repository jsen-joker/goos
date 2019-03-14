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
import org.jsen.goos.client.listener.Listener;

/**
 * Config Service Interface
 *
 * @author Nacos
 */
public interface ConfigService {

    /**
     * Get config
     *
     * @param dataId    dataId
     * @param group     group
     * @param timeoutMs read timeout
     * @return config value
     * @throws GoosClientException GoosClientException
     */
    String getConfig(String dataId, String group, long timeoutMs) throws GoosClientException;

    /**
     * Add a listener to the configuration, after the server modified the
     * configuration, the client will use the incoming listener callback.
     * Recommended asynchronous processing, the application can implement the
     * getExecutor method in the ManagerListener, provide a thread pool of
     * execution. If provided, use the main thread callback, May block other
     * configurations or be blocked by other configurations.
     *
     * @param dataId   dataId
     * @param group    group
     * @param listener listener
     * @throws GoosClientException GoosClientException
     */
    void addListener(String dataId, String group, Listener listener) throws GoosClientException;

    /**
     * Publish config.
     *
     * @param dataId  dataId
     * @param group   group
     * @param content content
     * @return Whether publish
     * @throws GoosClientException GoosClientException
     */
    boolean publishConfig(String dataId, String group, String content) throws GoosClientException;

    /**
     * Remove config
     *
     * @param dataId dataId
     * @param group  group
     * @return whether remove
     * @throws GoosClientException GoosClientException
     */
    boolean removeConfig(String dataId, String group) throws GoosClientException;

    /**
     * Remove listener
     *
     * @param dataId   dataId
     * @param group    group
     * @param listener listener
     */
    void removeListener(String dataId, String group, Listener listener);

    /**
     * Get server status
     *
     * @return whether health
     */
    String getServerStatus();

}
