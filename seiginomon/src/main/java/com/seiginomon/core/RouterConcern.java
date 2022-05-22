package com.seiginomon.core;

import java.util.List;

import org.springframework.cloud.gateway.filter.ratelimit.RedisRateLimiter;
import org.springframework.cloud.gateway.route.RouteLocator;
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.WebSession;

import lombok.Value;
import reactor.core.publisher.Mono;

@Configuration
@RestController
public class RouterConcern {

	// Rate limit configuration -----------------
	@Bean
	public RedisRateLimiter redisRateLimiter() {
		return new RedisRateLimiter(5, 7); // replenishRate, burstCapacity
	}
	// ------------------------------------------

	// Gateway/Edge configuration --------------------
	@Bean
	public RouteLocator appRouteLocator(RouteLocatorBuilder builder, RedisRateLimiter redisRateLimiter) {
		//@formatter:off
		return builder.routes().route(r -> r.path("/dogs")
		    .filters(f -> f
	    		.rewritePath("/dogs", "/api/facts")
	    		.requestRateLimiter(rl -> rl.setRateLimiter(redisRateLimiter))
	    		.circuitBreaker(cb -> cb
    				.setName("dogsApiCircuitBreaker")
    				.setFallbackUri("forward:/dogs_api_fallback")
    			)
	    	).uri("https://dog-api.kinduff.com")
		).build();
		//@formatter:on
	}
	// ------------------------------------------

	// Fallback example -------------------------
	@Value // Upgrade and convert into record
	private static class DogFact {
		private final List<String> facts;
		private final Boolean success;
		private String session;
	}

	@GetMapping("/dogs_api_fallback")
	public Mono<DogFact> dogsApiFallback(WebSession session) {
		//@formatter:off
		return Mono.just(
			new DogFact(
				List.of("A dog could detect a teaspoon of sugar if you added it to an Olympic-sized swimming pool full of water."),
				true,
				session.getId()
			)
		);
		//@formatter:on
	}
	// ------------------------------------------
}
