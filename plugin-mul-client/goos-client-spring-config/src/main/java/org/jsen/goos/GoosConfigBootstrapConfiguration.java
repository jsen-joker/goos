package org.jsen.goos;

import org.jsen.goos.locator.GoosPropertySourceLocator;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * <p>
 * </p>
 *
 * @author jsen
 * @since 2019-03-14
 */
@Configuration
public class GoosConfigBootstrapConfiguration {


    @Bean
    public GoosPropertySourceLocator goosPropertySourceLocator() {
        return new GoosPropertySourceLocator();
    }



    @Bean
    @ConditionalOnMissingBean
    public GoosConfigProperties goosConfigProperties() {
        return new GoosConfigProperties();
    }

}
