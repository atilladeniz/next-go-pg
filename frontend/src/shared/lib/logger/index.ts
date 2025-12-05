import pino from "pino"

const isDevelopment = process.env.NODE_ENV === "development"
const isServer = typeof window === "undefined"

// Create base logger configuration
const baseConfig: pino.LoggerOptions = {
	level: process.env.LOG_LEVEL || (isDevelopment ? "debug" : "info"),
	// Redact sensitive fields
	redact: {
		paths: ["password", "token", "authorization", "cookie", "*.password", "*.token"],
		censor: "[REDACTED]",
	},
}

// Server-side logger with pretty printing in development
const createServerLogger = (): pino.Logger => {
	if (isDevelopment) {
		// Dynamic import for pino-pretty in development
		return pino({
			...baseConfig,
			transport: {
				target: "pino-pretty",
				options: {
					colorize: true,
					translateTime: "SYS:standard",
					ignore: "pid,hostname",
				},
			},
		})
	}

	// Production: JSON logs
	return pino({
		...baseConfig,
		formatters: {
			level: (label) => ({ level: label }),
		},
		timestamp: pino.stdTimeFunctions.isoTime,
	})
}

// Client-side logger (browser)
const createClientLogger = (): pino.Logger => {
	return pino({
		...baseConfig,
		browser: {
			asObject: true,
			// In production, you could send logs to a service
			// transmit: {
			//   send: (level, logEvent) => {
			//     // Send to PostHog, Sentry, etc.
			//   }
			// }
		},
	})
}

// Export singleton logger
export const logger = isServer ? createServerLogger() : createClientLogger()

// Convenience methods with context
export const log = {
	debug: (msg: string, data?: Record<string, unknown>) => logger.debug(data, msg),
	info: (msg: string, data?: Record<string, unknown>) => logger.info(data, msg),
	warn: (msg: string, data?: Record<string, unknown>) => logger.warn(data, msg),
	error: (msg: string, data?: Record<string, unknown>) => logger.error(data, msg),

	// HTTP request logging
	request: (method: string, path: string, statusCode: number, duration: number) => {
		const level = statusCode >= 500 ? "error" : statusCode >= 400 ? "warn" : "info"
		logger[level]({ method, path, statusCode, duration }, "HTTP request")
	},

	// Auth events
	authSuccess: (userId: string) => {
		logger.info({ userId }, "Authentication successful")
	},

	authFailure: (reason: string) => {
		logger.warn({ reason }, "Authentication failed")
	},

	// API call logging
	apiCall: (endpoint: string, method: string, duration: number, success: boolean) => {
		const level = success ? "debug" : "error"
		logger[level]({ endpoint, method, duration, success }, "API call")
	},

	// Business events (for PostHog/analytics integration)
	event: (eventName: string, properties?: Record<string, unknown>) => {
		logger.info({ eventName, ...properties }, "Business event")
	},
}

// Child logger factory for component/feature-specific logging
export const createLogger = (component: string) => {
	const child = logger.child({ component })

	return {
		debug: (msg: string, data?: Record<string, unknown>) => child.debug(data, msg),
		info: (msg: string, data?: Record<string, unknown>) => child.info(data, msg),
		warn: (msg: string, data?: Record<string, unknown>) => child.warn(data, msg),
		error: (msg: string, data?: Record<string, unknown>) => child.error(data, msg),
	}
}

// Type exports
export type Logger = typeof log
export type ComponentLogger = ReturnType<typeof createLogger>
