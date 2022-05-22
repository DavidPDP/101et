package com.seiginomon.core;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;
import org.springframework.session.data.redis.config.annotation.web.server.EnableRedisWebSession;

@Configuration
@EnableRedisWebSession
public class SecurityConcern {
	
	// Security filters configuration -----------
	@Bean
	public SecurityWebFilterChain securityWebFilterChain(ServerHttpSecurity http) {
		http.formLogin();
		http.authorizeExchange().anyExchange().authenticated();
		return http.build();
	}
	// -------------------------------------------
	
}
