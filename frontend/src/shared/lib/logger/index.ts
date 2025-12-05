import pino from "pino"
import { sendToLoki } from "./loki-transport"

const isDevelopment = process.env.NODE_ENV === "development"
const isServer = typeof window === "undefined"

// ============================================================================
// Log Categories (for filtering)
// ============================================================================

export const LogCategory = {
	HTTP: "http",
	AUTH: "auth",
	API: "api",
	UI: "ui",
	BUSINESS: "business",
	PERFORMANCE: "performance",
	ERROR: "error",
} as const

export type LogCategoryType = (typeof LogCategory)[keyof typeof LogCategory]

// ============================================================================
// Logger Configuration
// ============================================================================

const baseConfig: pino.LoggerOptions = {
	level: process.env.LOG_LEVEL || (isDevelopment ? "debug" : "info"),
	// Redact sensitive fields
	redact: {
		paths: [
			"password",
			"token",
			"authorization",
			"cookie",
			"*.password",
			"*.token",
			"*.secret",
			"creditCard",
			"*.creditCard",
		],
		censor: "[REDACTED]",
	},
}

// Server-side logger with pretty printing in development
const createServerLogger = (): pino.Logger => {
	if (isDevelopment) {
		return pino({
			...baseConfig,
			transport: {
				target: "pino-pretty",
				options: {
					colorize: true,
					translateTime: "HH:MM:ss",
					ignore: "pid,hostname",
					messageFormat: "{category} | {msg}",
				},
			},
		})
	}

	// Production: JSON logs with caller info
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
			// Custom serializers for browser
			serialize: true,
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

// ============================================================================
// User Context (for adding user info to logs)
// ============================================================================

interface UserContext {
	userId?: string
	userName?: string
	sessionId?: string
}

let currentUserContext: UserContext = {}

// Set user context for all subsequent logs
export function setUserContext(context: UserContext) {
	currentUserContext = context
}

// Clear user context (on logout)
export function clearUserContext() {
	currentUserContext = {}
}

// Get logger with current user context
function getContextualLogger() {
	if (Object.keys(currentUserContext).length === 0) {
		return logger
	}
	return logger.child(currentUserContext)
}

// ============================================================================
// Structured Logging Helpers
// ============================================================================

// Helper to log to both Pino and Loki
function logWithLoki(
	level: "debug" | "info" | "warn" | "error",
	msg: string,
	data?: Record<string, unknown>,
) {
	const contextLogger = getContextualLogger()
	const logData = { ...currentUserContext, ...data }
	contextLogger[level](logData, msg)

	// Also send to Loki (server-side only)
	if (isServer) {
		sendToLoki(level, msg, logData)
	}
}

export const log = {
	// Basic logging with category
	debug: (msg: string, data?: Record<string, unknown>) => {
		logWithLoki("debug", msg, data)
	},

	info: (msg: string, data?: Record<string, unknown>) => {
		logWithLoki("info", msg, data)
	},

	warn: (msg: string, data?: Record<string, unknown>) => {
		logWithLoki("warn", msg, data)
	},

	error: (msg: string, data?: Record<string, unknown>) => {
		logWithLoki("error", msg, data)
	},

	// ========================================================================
	// HTTP/API Logging
	// ========================================================================

	// HTTP request logging (for API calls)
	request: (
		method: string,
		path: string,
		statusCode: number,
		duration: number,
		extra?: Record<string, unknown>,
	) => {
		const level = statusCode >= 500 ? "error" : statusCode >= 400 ? "warn" : "info"
		const data = {
			category: LogCategory.HTTP,
			method,
			path,
			status: statusCode,
			duration_ms: duration,
			...extra,
		}
		logWithLoki(level, "HTTP request", data)
	},

	// API call tracking (for Orval/TanStack Query)
	apiCall: (
		endpoint: string,
		method: string,
		duration: number,
		success: boolean,
		extra?: Record<string, unknown>,
	) => {
		const level = success ? "debug" : "error"
		logWithLoki(level, "API call", {
			category: LogCategory.API,
			endpoint,
			method,
			duration_ms: duration,
			success,
			...extra,
		})
	},

	// ========================================================================
	// Authentication Logging
	// ========================================================================

	authSuccess: (action: string, userId: string, extra?: Record<string, unknown>) => {
		logWithLoki("info", "Auth event", {
			category: LogCategory.AUTH,
			action,
			user_id: userId,
			success: true,
			...extra,
		})
	},

	authFailure: (action: string, reason: string, extra?: Record<string, unknown>) => {
		logWithLoki("warn", "Auth event", {
			category: LogCategory.AUTH,
			action,
			reason,
			success: false,
			...extra,
		})
	},

	// ========================================================================
	// Business Event Logging (for PostHog/analytics)
	// ========================================================================

	event: (eventName: string, properties?: Record<string, unknown>) => {
		logWithLoki("info", "Business event", {
			category: LogCategory.BUSINESS,
			event_name: eventName,
			...properties,
		})
	},

	// ========================================================================
	// UI/Component Logging
	// ========================================================================

	component: (componentName: string, action: string, data?: Record<string, unknown>) => {
		logWithLoki("debug", "Component event", {
			category: LogCategory.UI,
			component: componentName,
			action,
			...data,
		})
	},

	// ========================================================================
	// Performance Logging
	// ========================================================================

	// Log slow operations
	slowOperation: (operation: string, duration: number, threshold: number) => {
		if (duration > threshold) {
			logWithLoki("warn", "Slow operation", {
				category: LogCategory.PERFORMANCE,
				operation,
				duration_ms: duration,
				threshold_ms: threshold,
				slow: true,
			})
		}
	},

	// Navigation timing
	navigation: (route: string, duration: number) => {
		logWithLoki("info", "Navigation", {
			category: LogCategory.PERFORMANCE,
			route,
			duration_ms: duration,
		})
	},

	// ========================================================================
	// Error Logging (enhanced)
	// ========================================================================

	// Log error with stack trace
	exception: (error: Error, context?: Record<string, unknown>) => {
		logWithLoki("error", "Exception", {
			category: LogCategory.ERROR,
			error_name: error.name,
			error_message: error.message,
			stack: isDevelopment ? error.stack : undefined,
			...context,
		})
	},

	// Log unhandled errors
	unhandled: (error: unknown, source: string) => {
		const errorObj = error instanceof Error ? error : new Error(String(error))
		logWithLoki("error", "Unhandled error", {
			category: LogCategory.ERROR,
			source,
			error_name: errorObj.name,
			error_message: errorObj.message,
			stack: isDevelopment ? errorObj.stack : undefined,
		})
	},
}

// ============================================================================
// Child Logger Factory (for component/feature-specific logging)
// ============================================================================

export interface ComponentLogger {
	debug: (msg: string, data?: Record<string, unknown>) => void
	info: (msg: string, data?: Record<string, unknown>) => void
	warn: (msg: string, data?: Record<string, unknown>) => void
	error: (msg: string, data?: Record<string, unknown>) => void
	timed: (operation: string) => () => void
}

export const createLogger = (component: string): ComponentLogger => {
	const child = logger.child({ component })

	return {
		debug: (msg: string, data?: Record<string, unknown>) => {
			child.debug({ ...currentUserContext, ...data }, msg)
		},
		info: (msg: string, data?: Record<string, unknown>) => {
			child.info({ ...currentUserContext, ...data }, msg)
		},
		warn: (msg: string, data?: Record<string, unknown>) => {
			child.warn({ ...currentUserContext, ...data }, msg)
		},
		error: (msg: string, data?: Record<string, unknown>) => {
			child.error({ ...currentUserContext, ...data }, msg)
		},
		// Timed operation helper
		timed: (operation: string) => {
			const start = Date.now()
			return () => {
				const duration = Date.now() - start
				child.debug(
					{ ...currentUserContext, operation, duration_ms: duration },
					"Operation completed",
				)
			}
		},
	}
}

// ============================================================================
// Performance Helpers
// ============================================================================

// Measure async operation duration
export async function withTiming<T>(
	operation: string,
	fn: () => Promise<T>,
	threshold?: number,
): Promise<T> {
	const start = Date.now()
	try {
		return await fn()
	} finally {
		const duration = Date.now() - start
		if (threshold && duration > threshold) {
			log.slowOperation(operation, duration, threshold)
		} else {
			log.debug(`${operation} completed`, { duration_ms: duration })
		}
	}
}

// ============================================================================
// Type Exports
// ============================================================================

export type Logger = typeof log
