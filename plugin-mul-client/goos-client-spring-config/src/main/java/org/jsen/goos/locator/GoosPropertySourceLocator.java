package org.jsen.goos.locator;

import org.jsen.goos.GoosConfigProperties;
import org.jsen.goos.client.ConfigService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cloud.bootstrap.config.PropertySourceLocator;
import org.springframework.core.annotation.Order;
import org.springframework.core.env.CompositePropertySource;
import org.springframework.core.env.Environment;
import org.springframework.core.env.PropertySource;
import org.springframework.util.StringUtils;

/**
 * <p>
 * </p>
 *
 * @author jsen
 * @since 2019-03-14
 */
@Order(0)
public class GoosPropertySourceLocator implements PropertySourceLocator {
    private static final Logger logger = LoggerFactory
            .getLogger(GoosPropertySourceLocator.class);
    private static final String NACOS_PROPERTY_SOURCE_NAME = "GOOS";
    private static final String SEP1 = "-";
    private static final String DOT = ".";

    @Autowired
    private GoosConfigProperties goosConfigProperties;

    private GoosPropertySourceBuilder goosPropertySourceBuilder;

    @Override
    public PropertySource<?> locate(Environment env) {

        ConfigService configService = goosConfigProperties.configServiceInstance();

        if (null == configService) {
            logger.warn(
                    "no instance of config service found, can't load config from nacos");
            return null;
        }
        long timeout = goosConfigProperties.getTimeout();
        goosPropertySourceBuilder = new GoosPropertySourceBuilder(configService,
                timeout);

        String name = goosConfigProperties.getName();

        String nacosGroup = goosConfigProperties.getGroup();
        String dataIdPrefix = goosConfigProperties.getPrefix();
        if (StringUtils.isEmpty(dataIdPrefix)) {
            dataIdPrefix = name;
        }

        String fileExtension = goosConfigProperties.getFileExtension();

        CompositePropertySource composite = new CompositePropertySource(
                NACOS_PROPERTY_SOURCE_NAME);

        loadApplicationConfiguration(composite, nacosGroup, dataIdPrefix, fileExtension);

        return composite;
    }

    private void loadApplicationConfiguration(
            CompositePropertySource compositePropertySource, String nacosGroup,
            String dataIdPrefix, String fileExtension) {
        loadNacosDataIfPresent(compositePropertySource,
                dataIdPrefix + DOT + fileExtension, nacosGroup, fileExtension);
        for (String profile : goosConfigProperties.getActiveProfiles()) {
            String dataId = dataIdPrefix + SEP1 + profile + DOT + fileExtension;
            loadNacosDataIfPresent(compositePropertySource, dataId, nacosGroup,
                    fileExtension);
        }
    }

    private void loadNacosDataIfPresent(final CompositePropertySource composite,
                                        final String dataId, final String group, String fileExtension) {
        GoosPropertySource ps = goosPropertySourceBuilder.build(dataId, group,
                fileExtension);
        if (ps != null) {
            composite.addFirstPropertySource(ps);
        }
    }
}
