<page>
  <title>Introduction | Better Auth</title>
  <url>https://www.better-auth.com/docs</url>
  <content>Better Auth is a framework-agnostic, universal authentication and authorization framework for TypeScript. It provides a comprehensive set of features out of the box and includes a plugin ecosystem that simplifies adding advanced functionalities. Whether you need 2FA, passkey, multi-tenancy, multi-session support, or even enterprise features like SSO, creating your own IDP, it lets you focus on building your application instead of reinventing the wheel.

Better Auth aims to be the most comprehensive auth library. It provides a wide range of features out of the box and allows you to extend it with plugins. Here are some of the features:

...and much more!

* * *

### [LLMs.txt](#llmstxt)

Better Auth exposes an `LLMs.txt` that helps AI models understand how to integrate and interact with your authentication system. See it at [https://better-auth.com/llms.txt](https://better-auth.com/llms.txt).

### [MCP](#mcp)

Better Auth provides an MCP server so you can use it with any AI model that supports the Model Context Protocol (MCP).

#### [CLI Options](#cli-options)

Use the Better Auth CLI to easily add the MCP server to your preferred client:

#### [Manual Configuration](#manual-configuration)

Alternatively, you can manually configure the MCP server for each client:

We provide a first‑party MCP, powered by [Chonkie](https://chonkie.ai/). You can alternatively use [`context7`](https://context7.com/) and other MCP providers.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/introduction.mdx)</content>
</page>

<page>
  <title>Comparison | Better Auth</title>
  <url>https://www.better-auth.com/docs/comparison</url>
  <content>> Comparison is the thief of joy.

Here are non detailed reasons why you may want to use Better Auth over other auth libraries and services.

### [vs Other Auth Libraries](#vs-other-auth-libraries)

*   **Framework agnostic** - Works with any framework, not just specific ones
*   **Advanced features built-in** - 2FA, multi-tenancy, multi-session, rate limiting, and many more
*   **Plugin system** - Extend functionality without forking or complex workarounds
*   **Full control** - Customize auth flows exactly how you want

### [vs Self-Hosted Auth Servers](#vs-self-hosted-auth-servers)

*   **No separate infrastructure** - Runs in your app, users stay in your database
*   **Zero server maintenance** - No auth servers to deploy, monitor, or update
*   **Complete feature set** - Everything you need without the operational overhead

### [vs Managed Auth Services](#vs-managed-auth-services)

*   **Keep your data** - Users stay in your database, not a third-party service
*   **No per-user costs** - Scale without worrying about auth billing
*   **Single source of truth** - All user data in one place

### [vs Rolling Your Own](#vs-rolling-your-own)

*   **Security handled** - Battle-tested auth flows and security practices
*   **Focus on your product** - Spend time on features that matter to your business
*   **Plugin extensibility** - Add custom features without starting from scratch

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/comparison.mdx)</content>
</page>

<page>
  <title>Installation | Better Auth</title>
  <url>https://www.better-auth.com/docs/installation</url>
  <content>### [Install the Package](#install-the-package)

Let's start by adding Better Auth to your project:

If you're using a separate client and server setup, make sure to install Better Auth in both parts of your project.

### [Set Environment Variables](#set-environment-variables)

Create a `.env` file in the root of your project and add the following environment variables:

1.  **Secret Key**

A secret value used for encryption and hashing. It must be at least 32 characters and generated with high entropy. Click the button below to generate one. You can also use `openssl rand -base64 32` to generate one.

.env

    BETTER_AUTH_SECRET=

2.  **Set Base URL**

.env

    BETTER_AUTH_URL=http://localhost:3000 # Base URL of your app

### [Create A Better Auth Instance](#create-a-better-auth-instance)

Create a file named `auth.ts` in one of these locations:

*   Project root
*   `lib/` folder
*   `utils/` folder

You can also nest any of these folders under `src/`, `app/` or `server/` folder. (e.g. `src/lib/auth.ts`, `app/lib/auth.ts`).

And in this file, import Better Auth and create your auth instance. Make sure to export the auth instance with the variable name `auth` or as a `default` export.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      //...
    });

### [Configure Database](#configure-database)

Better Auth requires a database to store user data. You can easily configure Better Auth to use SQLite, PostgreSQL, or MySQL, and more!

You can also configure Better Auth to work in a stateless mode if you don't configure a database. See [Stateless Session Management](https://www.better-auth.com/docs/concepts/session-management#stateless-session-management) for more information. Note that most plugins will require a database.

Alternatively, if you prefer to use an ORM, you can use one of the built-in adapters.

If your database is not listed above, check out our other supported [databases](https://www.better-auth.com/docs/adapters/other-relational-databases) for more information, or use one of the supported ORMs.

### [Create Database Tables](#create-database-tables)

Better Auth includes a CLI tool to help manage the schema required by the library.

*   **Generate**: This command generates an ORM schema or SQL migration file.

If you're using Kysely, you can apply the migration directly with `migrate` command below. Use `generate` only if you plan to apply the migration manually.

Terminal

    npx @better-auth/cli generate

*   **Migrate**: This command creates the required tables directly in the database. (Available only for the built-in Kysely adapter)

Terminal

    npx @better-auth/cli migrate

see the [CLI documentation](https://www.better-auth.com/docs/concepts/cli) for more information.

If you instead want to create the schema manually, you can find the core schema required in the [database section](https://www.better-auth.com/docs/concepts/database#core-schema).

### [Authentication Methods](#authentication-methods)

Configure the authentication methods you want to use. Better Auth comes with built-in support for email/password, and social sign-on providers.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      //...other options
      emailAndPassword: { 
        enabled: true, 
      }, 
      socialProviders: { 
        github: { 
          clientId: process.env.GITHUB_CLIENT_ID as string, 
          clientSecret: process.env.GITHUB_CLIENT_SECRET as string, 
        }, 
      }, 
    });

### [Mount Handler](#mount-handler)

To handle API requests, you need to set up a route handler on your server.

Create a new file or route in your framework's designated catch-all route handler. This route should handle requests for the path `/api/auth/*` (unless you've configured a different base path).

Better Auth supports any backend framework with standard Request and Response objects and offers helper functions for popular frameworks.

### [Create Client Instance](#create-client-instance)

The client-side library helps you interact with the auth server. Better Auth comes with a client for all the popular web frameworks, including vanilla JavaScript.

1.  Import `createAuthClient` from the package for your framework (e.g., "better-auth/react" for React).
2.  Call the function to create your client.
3.  Pass the base URL of your auth server. (If the auth server is running on the same domain as your client, you can skip this step.)

If you're using a different base path other than `/api/auth` make sure to pass the whole URL including the path. (e.g. `http://localhost:3000/custom-path/auth`)

Tip: You can also export specific methods if you prefer:

    export const { signIn, signUp, useSession } = createAuthClient()

### [🎉 That's it!](#-thats-it)

That's it! You're now ready to use better-auth in your application. Continue to [basic usage](https://www.better-auth.com/docs/basic-usage) to learn how to use the auth instance to sign in users.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/installation.mdx)</content>
</page>

<page>
  <title>Introduction | Better Auth</title>
  <url>https://www.better-auth.com/docs/introduction</url>
  <content>Better Auth is a framework-agnostic, universal authentication and authorization framework for TypeScript. It provides a comprehensive set of features out of the box and includes a plugin ecosystem that simplifies adding advanced functionalities. Whether you need 2FA, passkey, multi-tenancy, multi-session support, or even enterprise features like SSO, creating your own IDP, it lets you focus on building your application instead of reinventing the wheel.

Better Auth aims to be the most comprehensive auth library. It provides a wide range of features out of the box and allows you to extend it with plugins. Here are some of the features:

...and much more!

* * *

### [LLMs.txt](#llmstxt)

Better Auth exposes an `LLMs.txt` that helps AI models understand how to integrate and interact with your authentication system. See it at [https://better-auth.com/llms.txt](https://better-auth.com/llms.txt).

### [MCP](#mcp)

Better Auth provides an MCP server so you can use it with any AI model that supports the Model Context Protocol (MCP).

#### [CLI Options](#cli-options)

Use the Better Auth CLI to easily add the MCP server to your preferred client:

#### [Manual Configuration](#manual-configuration)

Alternatively, you can manually configure the MCP server for each client:

We provide a first‑party MCP, powered by [Chonkie](https://chonkie.ai/). You can alternatively use [`context7`](https://context7.com/) and other MCP providers.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/introduction.mdx)</content>
</page>

<page>
  <title>Basic Usage | Better Auth</title>
  <url>https://www.better-auth.com/docs/basic-usage</url>
  <content>Better Auth provides built-in authentication support for:

*   **Email and password**
*   **Social provider (Google, GitHub, Apple, and more)**

But also can easily be extended using plugins, such as: [username](https://www.better-auth.com/docs/plugins/username), [magic link](https://www.better-auth.com/docs/plugins/magic-link), [passkey](https://www.better-auth.com/docs/plugins/passkey), [email-otp](https://www.better-auth.com/docs/plugins/email-otp), and more.

To enable email and password authentication:

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        emailAndPassword: {    
            enabled: true
        } 
    })

### [Sign Up](#sign-up)

To sign up a user you need to call the client method `signUp.email` with the user's information.

sign-up.ts

    import { authClient } from "@/lib/auth-client"; //import the auth client
    
    const { data, error } = await authClient.signUp.email({
            email, // user email address
            password, // user password -> min 8 characters by default
            name, // user display name
            image, // User image URL (optional)
            callbackURL: "/dashboard" // A URL to redirect to after the user verifies their email (optional)
        }, {
            onRequest: (ctx) => {
                //show loading
            },
            onSuccess: (ctx) => {
                //redirect to the dashboard or sign in page
            },
            onError: (ctx) => {
                // display the error message
                alert(ctx.error.message);
            },
    });

By default, the users are automatically signed in after they successfully sign up. To disable this behavior you can set `autoSignIn` to `false`.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        emailAndPassword: {
        	enabled: true,
        	autoSignIn: false //defaults to true
      },
    })

### [Sign In](#sign-in)

To sign a user in, you can use the `signIn.email` function provided by the client.

sign-in

    const { data, error } = await authClient.signIn.email({
            /**
             * The user email
             */
            email,
            /**
             * The user password
             */
            password,
            /**
             * A URL to redirect to after the user verifies their email (optional)
             */
            callbackURL: "/dashboard",
            /**
             * remember the user session after the browser is closed. 
             * @default true
             */
            rememberMe: false
    }, {
        //callbacks
    })

Always invoke client methods from the client side. Don't call them from the server.

### [Server-Side Authentication](#server-side-authentication)

To authenticate a user on the server, you can use the `auth.api` methods.

server.ts

    import { auth } from "./auth"; // path to your Better Auth server instance
    
    const response = await auth.api.signInEmail({
        body: {
            email,
            password
        },
        asResponse: true // returns a response object instead of data
    });

If the server cannot return a response object, you'll need to manually parse and set cookies. But for frameworks like Next.js we provide [a plugin](https://www.better-auth.com/docs/integrations/next#server-action-cookies) to handle this automatically

Better Auth supports multiple social providers, including Google, GitHub, Apple, Discord, and more. To use a social provider, you need to configure the ones you need in the `socialProviders` option on your `auth` object.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        socialProviders: { 
            github: { 
                clientId: process.env.GITHUB_CLIENT_ID!, 
                clientSecret: process.env.GITHUB_CLIENT_SECRET!, 
            } 
        }, 
    })

### [Sign in with social providers](#sign-in-with-social-providers)

To sign in using a social provider you need to call `signIn.social`. It takes an object with the following properties:

sign-in.ts

    import { authClient } from "@/lib/auth-client"; //import the auth client
    
    await authClient.signIn.social({
        /**
         * The social provider ID
         * @example "github", "google", "apple"
         */
        provider: "github",
        /**
         * A URL to redirect after the user authenticates with the provider
         * @default "/"
         */
        callbackURL: "/dashboard", 
        /**
         * A URL to redirect if an error occurs during the sign in process
         */
        errorCallbackURL: "/error",
        /**
         * A URL to redirect if the user is newly registered
         */
        newUserCallbackURL: "/welcome",
        /**
         * disable the automatic redirect to the provider. 
         * @default false
         */
        disableRedirect: true,
    });

You can also authenticate using `idToken` or `accessToken` from the social provider instead of redirecting the user to the provider's site. See social providers documentation for more details.

To signout a user, you can use the `signOut` function provided by the client.

user-card.tsx

    await authClient.signOut();

you can pass `fetchOptions` to redirect onSuccess

user-card.tsx

    await authClient.signOut({
      fetchOptions: {
        onSuccess: () => {
          router.push("/login"); // redirect to login page
        },
      },
    });

Once a user is signed in, you'll want to access the user session. Better Auth allows you to easily access the session data from both the server and client sides.

### [Client Side](#client-side)

#### [Use Session](#use-session)

Better Auth provides a `useSession` hook to easily access session data on the client side. This hook is implemented using nanostore and has support for each supported framework and vanilla client, ensuring that any changes to the session (such as signing out) are immediately reflected in your UI.

#### [Get Session](#get-session)

If you prefer not to use the hook, you can use the `getSession` method provided by the client.

user.tsx

    import { authClient } from "@/lib/auth-client" // import the auth client
    
    const { data: session, error } = await authClient.getSession()

You can also use it with client-side data-fetching libraries like [TanStack Query](https://tanstack.com/query/latest).

### [Server Side](#server-side)

The server provides a `session` object that you can use to access the session data. It requires request headers object to be passed to the `getSession` method.

**Example: Using some popular frameworks**

One of the unique features of Better Auth is a plugins ecosystem. It allows you to add complex auth related functionality with small lines of code.

Below is an example of how to add two factor authentication using two factor plugin.

### [Server Configuration](#server-configuration)

To add a plugin, you need to import the plugin and pass it to the `plugins` option of the auth instance. For example, to add two factor authentication, you can use the following code:

auth.ts

    import { betterAuth } from "better-auth"
    import { twoFactor } from "better-auth/plugins"
    
    export const auth = betterAuth({
        //...rest of the options
        plugins: [ 
            twoFactor() 
        ] 
    })

now two factor related routes and method will be available on the server.

### [Migrate Database](#migrate-database)

After adding the plugin, you'll need to add the required tables to your database. You can do this by running the `migrate` command, or by using the `generate` command to create the schema and handle the migration manually.

generating the schema:

terminal

    npx @better-auth/cli generate

using the `migrate` command:

terminal

    npx @better-auth/cli migrate

If you prefer adding the schema manually, you can check the schema required on the [two factor plugin](https://www.better-auth.com/docs/plugins/2fa#schema) documentation.

### [Client Configuration](#client-configuration)

Once we're done with the server, we need to add the plugin to the client. To do this, you need to import the plugin and pass it to the `plugins` option of the auth client. For example, to add two factor authentication, you can use the following code:

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { twoFactorClient } from "better-auth/client/plugins"; 
    
    const authClient = createAuthClient({
        plugins: [ 
            twoFactorClient({ 
                twoFactorPage: "/two-factor" // the page to redirect if a user needs to verify 2nd factor
            }) 
        ] 
    })

now two factor related methods will be available on the client.

profile.ts

    import { authClient } from "./auth-client"
    
    const enableTwoFactor = async() => {
        const data = await authClient.twoFactor.enable({
            password // the user password is required
        }) // this will enable two factor
    }
    
    const disableTwoFactor = async() => {
        const data = await authClient.twoFactor.disable({
            password // the user password is required
        }) // this will disable two factor
    }
    
    const signInWith2Factor = async() => {
        const data = await authClient.signIn.email({
            //...
        })
        //if the user has two factor enabled, it will redirect to the two factor page
    }
    
    const verifyTOTP = async() => {
        const data = await authClient.twoFactor.verifyTOTP({
            code: "123456", // the code entered by the user 
            /**
             * If the device is trusted, the user won't
             * need to pass 2FA again on the same device
             */
            trustDevice: true
        })
    }</content>
</page>

<page>
  <title>Session Management | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/session-management#stateless-session-management</url>
  <content>Better Auth manages session using a traditional cookie-based session management. The session is stored in a cookie and is sent to the server on every request. The server then verifies the session and returns the user data if the session is valid.

The session table stores the session data. The session table has the following fields:

*   `id`: Unique identifier for the session.
*   `token`: The session token. Which is also used as the session cookie.
*   `userId`: The user ID of the user.
*   `expiresAt`: The expiration date of the session.
*   `ipAddress`: The IP address of the user.
*   `userAgent`: The user agent of the user. It stores the user agent header from the request.

The session expires after 7 days by default. But whenever the session is used and the `updateAge` is reached, the session expiration is updated to the current time plus the `expiresIn` value.

You can change both the `expiresIn` and `updateAge` values by passing the `session` object to the `auth` configuration.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        //... other config options
        session: {
            expiresIn: 60 * 60 * 24 * 7, // 7 days
            updateAge: 60 * 60 * 24 // 1 day (every 1 day the session expiration is updated)
        }
    })

### [Disable Session Refresh](#disable-session-refresh)

You can disable session refresh so that the session is not updated regardless of the `updateAge` option.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        //... other config options
        session: {
            disableSessionRefresh: true
        }
    })

Some endpoints in Better Auth require the session to be **fresh**. A session is considered fresh if its `createdAt` is within the `freshAge` limit. By default, the `freshAge` is set to **1 day** (60 \* 60 \* 24).

You can customize the `freshAge` value by passing a `session` object in the `auth` configuration:

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        //... other config options
        session: {
            freshAge: 60 * 5 // 5 minutes (the session is fresh if created within the last 5 minutes)
        }
    })

To **disable the freshness check**, set `freshAge` to `0`:

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        //... other config options
        session: {
            freshAge: 0 // Disable freshness check
        }
    })

Better Auth provides a set of functions to manage sessions.

### [Get Session](#get-session)

The `getSession` function retrieves the current active session.

    import { authClient } from "@/lib/client"
    
    const { data: session } = await authClient.getSession()

To learn how to customize the session response check the [Customizing Session Response](#customizing-session-response) section.

### [Use Session](#use-session)

The `useSession` action provides a reactive way to access the current session.

    import { authClient } from "@/lib/client"
    
    const { data: session } = authClient.useSession()

### [List Sessions](#list-sessions)

The `listSessions` function returns a list of sessions that are active for the user.

auth-client.ts

    import { authClient } from "@/lib/client"
    
    const sessions = await authClient.listSessions()

### [Revoke Session](#revoke-session)

When a user signs out of a device, the session is automatically ended. However, you can also end a session manually from any device the user is signed into.

To end a session, use the `revokeSession` function. Just pass the session token as a parameter.

auth-client.ts

    await authClient.revokeSession({
        token: "session-token"
    })

### [Revoke Other Sessions](#revoke-other-sessions)

To revoke all other sessions except the current session, you can use the `revokeOtherSessions` function.

auth-client.ts

    await authClient.revokeOtherSessions()

### [Revoke All Sessions](#revoke-all-sessions)

To revoke all sessions, you can use the `revokeSessions` function.

auth-client.ts

    await authClient.revokeSessions()

### [Revoking Sessions on Password Change](#revoking-sessions-on-password-change)

You can revoke all sessions when the user changes their password by passing `revokeOtherSessions` as true on `changePassword` function.

auth.ts

    await authClient.changePassword({
        newPassword: newPassword,
        currentPassword: currentPassword,
        revokeOtherSessions: true,
    })

### [Cookie Cache](#cookie-cache)

Calling your database every time `useSession` or `getSession` is invoked isn't ideal, especially if sessions don't change frequently. Cookie caching handles this by storing session data in a short-lived, signed cookie—similar to how JWT access tokens are used with refresh tokens.

When cookie caching is enabled, the server can check session validity from the cookie itself instead of hitting the database each time. The cookie is signed to prevent tampering, and a short `maxAge` ensures that the session data gets refreshed regularly. If a session is revoked or expires, the cookie will be invalidated automatically.

To turn on cookie caching, just set `session.cookieCache` in your auth config:

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        session: {
            cookieCache: {
                enabled: true,
                maxAge: 5 * 60 // Cache duration in seconds (5 minutes)
            }
        }
    });

#### [Cookie Cache Strategies](#cookie-cache-strategies)

Better Auth supports three different encoding strategies for cookie cache:

*   **`compact`** (default): Uses base64url encoding with HMAC-SHA256 signature. Most compact format with no JWT spec overhead. Best for performance and size.
*   **`jwt`**: Standard JWT with HMAC-SHA256 signature (HS256). Signed but not encrypted - readable by anyone but tamper-proof. Follows JWT spec for interoperability.
*   **`jwe`**: Uses JWE (JSON Web Encryption) with A256CBC-HS512 and HKDF key derivation. Fully encrypted tokens - neither readable nor tamperable. Most secure but largest size.

**Comparison:**

| Strategy | Size | Security | Readable | Interoperable | Use Case |
| --- | --- | --- | --- | --- | --- |
| `compact` | Smallest | Good (signed) | Yes | No | Performance-critical, internal use |
| `jwt` | Medium | Good (signed) | Yes | Yes | Need JWT compatibility, external integrations |
| `jwe` | Largest | Best (encrypted) | No | Yes | Sensitive data, maximum security |

auth.ts

    export const auth = betterAuth({
        session: {
            cookieCache: {
                enabled: true,
                maxAge: 5 * 60,
                strategy: "compact" // or "jwt" or "jwe"
            }
        }
    });

**Note:** All strategies are cryptographically secure and prevent tampering. The main differences are size, readability, and JWT spec compliance.

**When to use each:**

*   **Use `compact`** when you need maximum performance and smallest cookie size. Best for most applications where cookies are only used internally by Better Auth.
*   **Use `jwt`** when you need JWT compatibility for external systems, or when you want standard JWT tokens that can be verified by third-party tools. The tokens are readable (base64-encoded JSON) but tamper-proof.
*   **Use `jwe`** when you need maximum security and want to hide session data from the client. The tokens are fully encrypted and cannot be read without the secret key. Use this for sensitive data or compliance requirements.

If you want to disable returning from the cookie cache when fetching the session, you can pass `disableCookieCache:true` this will force the server to fetch the session from the database and also refresh the cookie cache.

auth-client.ts

    const session = await authClient.getSession({ query: {
        disableCookieCache: true
    }})

or on the server

server.ts

    await auth.api.getSession({
        query: {
            disableCookieCache: true,
        }, 
        headers: req.headers, // pass the headers
    });

Better Auth supports stateless session management without any database. This means that the session data is stored in a signed/encrypted cookie and the server never queries a database to validate sessions - it simply verifies the cookie signature and checks expiration.

### [Basic Stateless Setup](#basic-stateless-setup)

If you don't pass a database configuration, Better Auth will automatically enable stateless mode.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        // No database configuration
        socialProviders: {
            google: {
                clientId: process.env.GOOGLE_CLIENT_ID,
                clientSecret: process.env.GOOGLE_CLIENT_SECRET,
            },
        },
    });

To manually enable stateless mode, you need to configure `cookieCache` and `account` with the following options:

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        session: {
            cookieCache: {
                enabled: true,
                maxAge: 7 * 24 * 60 * 60, // 7 days cache duration
                strategy: "jwe", // can be "jwt" or "compact"
                refreshCache: true, // Enable stateless refresh
            },
        },
        account: {
            storeStateStrategy: "cookie",
            storeAccountCookie: true, // Store account data after OAuth flow in a cookie (useful for database-less flows)
        }
    });

If you don't provide a database, by default we provide the above configuration for you.

### [Understanding `refreshCache`](#understanding-refreshcache)

The `refreshCache` option controls automatic cookie refresh **before expiry** without querying any database:

*   **`false`** (default): No automatic refresh. When the cookie cache expires (reaches `maxAge`), it will attempt to fetch from the database if available.
*   **`true`**: Enable automatic refresh with default settings. Refreshes when 80% of `maxAge` is reached (20% time remaining).
*   **`object`**: Custom refresh configuration with `updateAge` property.

auth.ts

    export const auth = betterAuth({
        session: {
            cookieCache: {
                enabled: true,
                maxAge: 300, // 5 minutes
                refreshCache: {
                    updateAge: 60 // Refresh when 60 seconds remain before expiry
                }
            }
        }
    });

### [Versioning Stateless Sessions](#versioning-stateless-sessions)

One of the biggest drawbacks of stateless sessions is that you can't invalidate session easily. To solve this with better auth, if you would like to invalidate all sessions, you can change the version of the cookie cache and re-deploy your application.

auth.ts

    export const auth = betterAuth({
        session: {
            cookieCache: {
                version: "2", // Change the version to invalidate all sessions
            }
        }
    });

This will invalidate all sessions that don't match the new version.

### [Stateless with Secondary Storage](#stateless-with-secondary-storage)

You can combine stateless sessions with secondary storage (Redis, etc.) for the best of both worlds:

auth.ts

    import { betterAuth } from "better-auth"
    import { redis } from "./redis"
    
    export const auth = betterAuth({
        // No primary database needed
        secondaryStorage: {
            get: async (key) => await redis.get(key),
            set: async (key, value, ttl) => await redis.set(key, value, "EX", ttl),
            delete: async (key) => await redis.del(key)
        },
        session: {
            cookieCache: {
                maxAge: 5 * 60, // 5 minutes (short-lived cookie)
                refreshCache: false // Disable stateless refresh
            }
        }
    });

This setup:

*   Uses cookies for session validation (no DB queries)
*   Uses Redis for storing session data and refreshing the cookie cache before expiry
*   You can revoke sessions from the secondary storage and the cookie cache will be invalidated on refresh

When you call `getSession` or `useSession`, the session data is returned as a `user` and `session` object. You can customize this response using the `customSession` plugin.

auth.ts

    import { customSession } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [
            customSession(async ({ user, session }) => {
                const roles = findUserRoles(session.session.userId);
                return {
                    roles,
                    user: {
                        ...user,
                        newField: "newField",
                    },
                    session
                };
            }),
        ],
    });

This will add `roles` and `user.newField` to the session response.

**Infer on the Client**

auth-client.ts

    import { customSessionClient } from "better-auth/client/plugins";
    import type { auth } from "@/lib/auth"; // Import the auth instance as a type
    
    const authClient = createAuthClient({
        plugins: [customSessionClient<typeof auth>()],
    });
    
    const { data } = authClient.useSession();
    const { data: sessionData } = await authClient.getSession();
    // data.roles
    // data.user.newField

### [Caveats on Customizing Session Response](#caveats-on-customizing-session-response)

1.  The passed `session` object to the callback does not infer fields added by plugins.

However, as a workaround, you can pull up your auth options and pass it to the plugin to infer the fields.

    import { betterAuth, BetterAuthOptions } from "better-auth";
    
    const options = {
      //...config options
      plugins: [
        //...plugins 
      ]
    } satisfies BetterAuthOptions;
    
    export const auth = betterAuth({
        ...options,
        plugins: [
            ...(options.plugins ?? []),
            customSession(async ({ user, session }, ctx) => {
                // now both user and session will infer the fields added by plugins and your custom fields
                return {
                    user,
                    session
                }
            }, options), // pass options here
        ]
    })

2.  When your server and client code are in separate projects or repositories, and you cannot import the `auth` instance as a type reference, type inference for custom session fields will not work on the client side.
3.  Session caching, including secondary storage or cookie cache, does not include custom fields. Each time the session is fetched, your custom session function will be called.

**Mutating the list-device-sessions endpoint** The `/multi-session/list-device-sessions` endpoint from the [multi-session](https://www.better-auth.com/docs/plugins/multi-session) plugin is used to list the devices that the user is signed into.

You can mutate the response of this endpoint by passing the `shouldMutateListDeviceSessionsEndpoint` option to the `customSession` plugin.

By default, we do not mutate the response of this endpoint.

auth.ts

    import { customSession } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [
            customSession(async ({ user, session }, ctx) => {
                return {
                    user,
                    session
                }
            }, {}, { shouldMutateListDeviceSessionsEndpoint: true }), 
        ],
    });</content>
</page>

<page>
  <title>Other Relational Databases | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/other-relational-databases</url>
  <content>Better Auth supports a wide range of database dialects out of the box thanks to [Kysely](https://kysely.dev/).

Any dialect supported by Kysely can be utilized with Better Auth, including capabilities for generating and migrating database schemas through the [CLI](https://www.better-auth.com/docs/concepts/cli).

You can see the full list of supported Kysely dialects [here](https://kysely.dev/docs/dialects).</content>
</page>

<page>
  <title>CLI | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/cli</url>
  <content>Better Auth comes with a built-in CLI to help you manage the database schemas, initialize your project, generate a secret key for your application, and gather diagnostic information about your setup.

The `generate` command creates the schema required by Better Auth. If you're using a database adapter like Prisma or Drizzle, this command will generate the right schema for your ORM. If you're using the built-in Kysely adapter, it will generate an SQL file you can run directly on your database.

### [Options](#options)

*   `--output` - Where to save the generated schema. For Prisma, it will be saved in prisma/schema.prisma. For Drizzle, it goes to schema.ts in your project root. For Kysely, it's an SQL file saved as schema.sql in your project root.
*   `--config` - The path to your Better Auth config file. By default, the CLI will search for an auth.ts file in **./**, **./utils**, **./lib**, or any of these directories under the `src` directory.
*   `--yes` - Skip the confirmation prompt and generate the schema directly.

The migrate command applies the Better Auth schema directly to your database. This is available if you're using the built-in Kysely adapter. For other adapters, you'll need to apply the schema using your ORM's migration tool.

### [Options](#options-1)

*   `--config` - The path to your Better Auth config file. By default, the CLI will search for an auth.ts file in **./**, **./utils**, **./lib**, or any of these directories under the `src` directory.
*   `--yes` - Skip the confirmation prompt and apply the schema directly.

**Using PostgreSQL with a non-default schema?**

The migrate command automatically detects your configured `search_path` and creates tables in the correct schema. See the [PostgreSQL adapter documentation](https://www.better-auth.com/docs/adapters/postgresql#use-a-non-default-schema) for configuration details.

The `init` command allows you to initialize Better Auth in your project.

### [Options](#options-2)

*   `--name` - The name of your application. (defaults to the `name` property in your `package.json`).
*   `--framework` - The framework your codebase is using. Currently, the only supported framework is `Next.js`.
*   `--plugins` - The plugins you want to use. You can specify multiple plugins by separating them with a comma.
*   `--database` - The database you want to use. Currently, the only supported database is `SQLite`.
*   `--package-manager` - The package manager you want to use. Currently, the only supported package managers are `npm`, `pnpm`, `yarn`, `bun` (defaults to the manager you used to initialize the CLI).

The `info` command provides diagnostic information about your Better Auth setup and environment. Useful for debugging and sharing when seeking support.

### [Output](#output)

The command displays:

*   **System**: OS, CPU, memory, Node.js version
*   **Package Manager**: Detected manager and version
*   **Better Auth**: Version and configuration (sensitive data auto-redacted)
*   **Frameworks**: Detected frameworks (Next.js, React, Vue, etc.)
*   **Databases**: Database clients and ORMs (Prisma, Drizzle, etc.)

### [Options](#options-3)

*   `--config` - Path to your Better Auth config file
*   `--json` - Output as JSON for sharing or programmatic use

### [Examples](#examples)

    # Basic usage
    npx @better-auth/cli@latest info
    
    # Custom config path
    npx @better-auth/cli@latest info --config ./config/auth.ts
    
    # JSON output
    npx @better-auth/cli@latest info --json > auth-info.json

Sensitive data like secrets, API keys, and database URLs are automatically replaced with `[REDACTED]` for safe sharing.

The CLI also provides a way to generate a secret key for your Better Auth instance.

**Error: Cannot find module X**

If you see this error, it means the CLI can't resolve imported modules in your Better Auth config file. We are working on a fix for many of these issues, but in the meantime, you can try the following:

*   Remove any import aliases in your config file and use relative paths instead. After running the CLI, you can revert to using aliases.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/concepts/cli.mdx)</content>
</page>

<page>
  <title>Database | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/database#core-schema</url>
  <content>Better Auth connects to a database to store data. The database will be used to store data such as users, sessions, and more. Plugins can also define their own database tables to store data.

You can pass a database connection to Better Auth by passing a supported database instance in the database options. You can learn more about supported database adapters in the [Other relational databases](https://www.better-auth.com/docs/adapters/other-relational-databases) documentation.

Better Auth comes with a CLI tool to manage database migrations and generate schema.

### [Running Migrations](#running-migrations)

The cli checks your database and prompts you to add missing tables or update existing ones with new columns. This is only supported for the built-in Kysely adapter. For other adapters, you can use the `generate` command to create the schema and handle the migration through your ORM.

    npx @better-auth/cli migrate

For PostgreSQL users: The migrate command supports non-default schemas. It automatically detects your `search_path` configuration and creates tables in the correct schema. See [PostgreSQL adapter](https://www.better-auth.com/docs/adapters/postgresql#use-a-non-default-schema) for details.

### [Generating Schema](#generating-schema)

Better Auth also provides a `generate` command to generate the schema required by Better Auth. The `generate` command creates the schema required by Better Auth. If you're using a database adapter like Prisma or Drizzle, this command will generate the right schema for your ORM. If you're using the built-in Kysely adapter, it will generate an SQL file you can run directly on your database.

    npx @better-auth/cli generate

See the [CLI](https://www.better-auth.com/docs/concepts/cli) documentation for more information on the CLI.

If you prefer adding tables manually, you can do that as well. The core schema required by Better Auth is described below and you can find additional schema required by plugins in the plugin documentation.

Secondary storage in Better Auth allows you to use key-value stores for managing session data, rate limiting counters, etc. This can be useful when you want to offload the storage of this intensive records to a high performance storage or even RAM.

### [Implementation](#implementation)

To use secondary storage, implement the `SecondaryStorage` interface:

    interface SecondaryStorage {
      get: (key: string) => Promise<unknown>; 
      set: (key: string, value: string, ttl?: number) => Promise<void>;
      delete: (key: string) => Promise<void>;
    }

Then, provide your implementation to the `betterAuth` function:

    betterAuth({
      // ... other options
      secondaryStorage: {
        // Your implementation here
      },
    });

**Example: Redis Implementation**

Here's a basic example using Redis:

    import { createClient } from "redis";
    import { betterAuth } from "better-auth";
    
    const redis = createClient();
    await redis.connect();
    
    export const auth = betterAuth({
    	// ... other options
    	secondaryStorage: {
    		get: async (key) => {
    			return await redis.get(key);
    		},
    		set: async (key, value, ttl) => {
    			if (ttl) await redis.set(key, value, { EX: ttl });
    			// or for ioredis:
    			// if (ttl) await redis.set(key, value, 'EX', ttl)
    			else await redis.set(key, value);
    		},
    		delete: async (key) => {
    			await redis.del(key);
    		}
    	}
    });

This implementation allows Better Auth to use Redis for storing session data and rate limiting counters. You can also add prefixes to the keys names.

Better Auth requires the following tables to be present in the database. The types are in `typescript` format. You can use corresponding types in your database.

### [User](#user)

Table Name: `user`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each user |
| name | string | \- | User's chosen display name |
| email | string | \- | User's email address for communication and login |
| emailVerified | boolean | \- | Whether the user's email is verified |
| image | string |  | User's image url |
| createdAt | Date | \- | Timestamp of when the user account was created |
| updatedAt | Date | \- | Timestamp of the last update to the user's information |

### [Session](#session)

Table Name: `session`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each session |
| userId | string |  | The ID of the user |
| token | string | \- | The unique session token |
| expiresAt | Date | \- | The time when the session expires |
| ipAddress | string |  | The IP address of the device |
| userAgent | string |  | The user agent information of the device |
| createdAt | Date | \- | Timestamp of when the session was created |
| updatedAt | Date | \- | Timestamp of when the session was updated |

### [Account](#account)

Table Name: `account`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each account |
| userId | string |  | The ID of the user |
| accountId | string | \- | The ID of the account as provided by the SSO or equal to userId for credential accounts |
| providerId | string | \- | The ID of the provider |
| accessToken | string |  | The access token of the account. Returned by the provider |
| refreshToken | string |  | The refresh token of the account. Returned by the provider |
| accessTokenExpiresAt | Date |  | The time when the access token expires |
| refreshTokenExpiresAt | Date |  | The time when the refresh token expires |
| scope | string |  | The scope of the account. Returned by the provider |
| idToken | string |  | The ID token returned from the provider |
| password | string |  | The password of the account. Mainly used for email and password authentication |
| createdAt | Date | \- | Timestamp of when the account was created |
| updatedAt | Date | \- | Timestamp of when the account was updated |

### [Verification](#verification)

Table Name: `verification`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each verification |
| identifier | string | \- | The identifier for the verification request |
| value | string | \- | The value to be verified |
| expiresAt | Date | \- | The time when the verification request expires |
| createdAt | Date | \- | Timestamp of when the verification request was created |
| updatedAt | Date | \- | Timestamp of when the verification request was updated |

Better Auth allows you to customize the table names and column names for the core schema. You can also extend the core schema by adding additional fields to the user and session tables.

### [Custom Table Names](#custom-table-names)

You can customize the table names and column names for the core schema by using the `modelName` and `fields` properties in your auth config:

auth.ts

    export const auth = betterAuth({
      user: {
        modelName: "users",
        fields: {
          name: "full_name",
          email: "email_address",
        },
      },
      session: {
        modelName: "user_sessions",
        fields: {
          userId: "user_id",
        },
      },
    });

Type inference in your code will still use the original field names (e.g., `user.name`, not `user.full_name`).

To customize table names and column name for plugins, you can use the `schema` property in the plugin config:

auth.ts

    import { betterAuth } from "better-auth";
    import { twoFactor } from "better-auth/plugins";
    
    export const auth = betterAuth({
      plugins: [
        twoFactor({
          schema: {
            user: {
              fields: {
                twoFactorEnabled: "two_factor_enabled",
                secret: "two_factor_secret",
              },
            },
          },
        }),
      ],
    });

### [Extending Core Schema](#extending-core-schema)

Better Auth provides a type-safe way to extend the `user` and `session` schemas. You can add custom fields to your auth config, and the CLI will automatically update the database schema. These additional fields will be properly inferred in functions like `useSession`, `signUp.email`, and other endpoints that work with user or session objects.

To add custom fields, use the `additionalFields` property in the `user` or `session` object of your auth config. The `additionalFields` object uses field names as keys, with each value being a `FieldAttributes` object containing:

*   `type`: The data type of the field (e.g., "string", "number", "boolean").
*   `required`: A boolean indicating if the field is mandatory.
*   `defaultValue`: The default value for the field (note: this only applies in the JavaScript layer; in the database, the field will be optional).
*   `input`: This determines whether a value can be provided when creating a new record (default: `true`). If there are additional fields, like `role`, that should not be provided by the user during signup, you can set this to `false`.

Here's an example of how to extend the user schema with additional fields:

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      user: {
        additionalFields: {
          role: {
            type: "string",
            required: false,
            defaultValue: "user",
            input: false, // don't allow user to set role
          },
          lang: {
            type: "string",
            required: false,
            defaultValue: "en",
          },
        },
      },
    });

Now you can access the additional fields in your application logic.

    //on signup
    const res = await auth.api.signUpEmail({
      email: "test@example.com",
      password: "password",
      name: "John Doe",
      lang: "fr",
    });
    
    //user object
    res.user.role; // > "admin"
    res.user.lang; // > "fr"

See the [TypeScript](https://www.better-auth.com/docs/concepts/typescript#inferring-additional-fields-on-client) documentation for more information on how to infer additional fields on the client side.

If you're using social / OAuth providers, you may want to provide `mapProfileToUser` to map the profile data to the user object. So, you can populate additional fields from the provider's profile.

**Example: Mapping Profile to User For `firstName` and `lastName`**

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      socialProviders: {
        github: {
          clientId: "YOUR_GITHUB_CLIENT_ID",
          clientSecret: "YOUR_GITHUB_CLIENT_SECRET",
          mapProfileToUser: (profile) => {
            return {
              firstName: profile.name.split(" ")[0],
              lastName: profile.name.split(" ")[1],
            };
          },
        },
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          mapProfileToUser: (profile) => {
            return {
              firstName: profile.given_name,
              lastName: profile.family_name,
            };
          },
        },
      },
    });

### [ID Generation](#id-generation)

Better Auth by default will generate unique IDs for users, sessions, and other entities. You can customize ID generation behavior using the `advanced.database.generateId` option.

#### [Option 1: Let Database Generate IDs](#option-1-let-database-generate-ids)

Setting `generateId` to `false` allows your database handle all ID generation: (outside of `generateId` being `serial` and some cases of `generateId` being `uuid`)

auth.ts

    import { betterAuth } from "better-auth";
    import { db } from "./db";
    
    export const auth = betterAuth({
      database: db,
      advanced: {
        database: {
          generateId: false, // "serial" for auto-incrementing numeric IDs
        },
      },
    });

#### [Option 2: Custom ID Generation Function](#option-2-custom-id-generation-function)

Use a function to generate IDs. You can return `false` or `undefined` from the function to let the database generate the ID for specific models:

auth.ts

    import { betterAuth } from "better-auth";
    import { db } from "./db";
    
    export const auth = betterAuth({
      database: db,
      advanced: {
        database: {
          generateId: (options) => {
            // Let database auto-generate for specific models
            if (options.model === "user" || options.model === "users") {
              return false; // Let database generate ID
            }
            // Generate UUIDs for other tables
            return crypto.randomUUID();
          },
        },
      },
    });

**Important**: Returning `false` or `undefined` from the `generateId` function lets the database handle ID generation for that specific model. Setting `generateId: false` (without a function) disables ID generation for **all** tables.

#### [Option 3: Consistent Custom ID Generator](#option-3-consistent-custom-id-generator)

Generate the same type of ID for all tables:

auth.ts

    import { betterAuth } from "better-auth";
    import { db } from "./db";
    
    export const auth = betterAuth({
      database: db,
      advanced: {
        database: {
          generateId: () => crypto.randomUUID(),
        },
      },
    });

### [Numeric IDs](#numeric-ids)

If you prefer auto-incrementing numeric IDs, you can set the `advanced.database.generateId` option to `"serial"`. Doing this will disable Better-Auth from generating IDs for any table, and will assume your database will generate the numeric ID automatically.

When enabled, the Better-Auth CLI will generate or migrate the schema with the `id` field as a numeric type for your database with auto-incrementing attributes associated with it.

    import { betterAuth } from "better-auth";
    import { db } from "./db";
    
    export const auth = betterAuth({
      database: db,
      advanced: {
        database: {
          generateId: "serial",
        },
      },
    });

Better-Auth will continue to infer the type of the `id` field as a `string` for the database, but will automatically convert it to a numeric type when fetching or inserting data from the database.

It's likely when grabbing `id` values returned from Better-Auth that you'll receive a string version of a number, this is normal. It's also expected that all id values passed to Better-Auth (eg via an endpoint body) is expected to be a string.

### [UUIDs](#uuids)

If you prefer UUIDs for the `id` field, you can set the `advanced.database.generateId` option to `"uuid"`. By default, Better-Auth will generate UUIDs for the `id` field for all tables, except adapters that use `PostgreSQL` where we allow the database to generate the UUID automatically.

By enabling this option, the Better-Auth CLI will generate or migrate the schema with the `id` field as a UUID type for your database. If the `uuid` type is not supported, we will generate a normal `string` type for the `id` field.

### [Mixed ID Types](#mixed-id-types)

If you need different ID types across tables (e.g., integer IDs for users, UUID strings for sessions/accounts/verification), use a `generateId` callback function.

auth.ts

    import { betterAuth } from "better-auth";
    import { db } from "./db";
    
    export const auth = betterAuth({
      database: db,
      user: {
        modelName: "users", // PostgreSQL: id serial primary key
      },
      session: {
        modelName: "session", // PostgreSQL: id text primary key
      },
      advanced: {
        database: {
          // Do NOT set useNumberId - it's global and affects all tables
          generateId: (options) => {
            if (options.model === "user" || options.model === "users") {
              return false; // Let PostgreSQL serial generate it
            }
            return crypto.randomUUID(); // UUIDs for session, account, verification
          },
        },
      },
    });

This configuration allows you to:

*   Use database auto-increment (serial, auto\_increment, etc.) for the users table
*   Generate UUIDs for all other tables (session, account, verification)
*   Maintain compatibility with existing schemas that use different ID types

**Use Case**: This is particularly useful when migrating from other authentication providers (like Clerk) where you have existing users with integer IDs but want UUID strings for new tables.

### [Database Hooks](#database-hooks)

Database hooks allow you to define custom logic that can be executed during the lifecycle of core database operations in Better Auth. You can create hooks for the following models: **user**, **session**, and **account**.

Additional fields are supported, however full type inference for these fields isn't yet supported. Improved type support is planned.

There are two types of hooks you can define:

#### [1\. Before Hook](#1-before-hook)

*   **Purpose**: This hook is called before the respective entity (user, session, or account) is created, updated, or deleted.
*   **Behavior**: If the hook returns `false`, the operation will be aborted. And If it returns a data object, it'll replace the original payload.

#### [2\. After Hook](#2-after-hook)

*   **Purpose**: This hook is called after the respective entity is created or updated.
*   **Behavior**: You can perform additional actions or modifications after the entity has been successfully created or updated.

**Example Usage**

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      databaseHooks: {
        user: {
          create: {
            before: async (user, ctx) => {
              // Modify the user object before it is created
              return {
                data: {
                  // Ensure to return Better-Auth named fields, not the original field names in your database.
                  ...user,
                  firstName: user.name.split(" ")[0],
                  lastName: user.name.split(" ")[1],
                },
              };
            },
            after: async (user) => {
              //perform additional actions, like creating a stripe customer
            },
          },
          delete: {
            before: async (user, ctx) => {
              console.log(`User ${user.email} is being deleted`);
              if (user.email.includes("admin")) {
                return false; // Abort deletion
              }
              
              return true; // Allow deletion
            },
            after: async (user) => {
              console.log(`User ${user.email} has been deleted`);
            },
          },
        },
        session: {
          delete: {
            before: async (session, ctx) => {
              console.log(`Session ${session.token} is being deleted`);
              if (session.userId === "admin-user-id") {
                return false; // Abort deletion
              }
              return true; // Allow deletion
            },
            after: async (session) => {
              console.log(`Session ${session.token} has been deleted`);
            },
          },
        },
      },
    });

#### [Throwing Errors](#throwing-errors)

If you want to stop the database hook from proceeding, you can throw errors using the `APIError` class imported from `better-auth/api`.

auth.ts

    import { betterAuth } from "better-auth";
    import { APIError } from "better-auth/api";
    
    export const auth = betterAuth({
      databaseHooks: {
        user: {
          create: {
            before: async (user, ctx) => {
              if (user.isAgreedToTerms === false) {
                // Your special condition.
                // Send the API error.
                throw new APIError("BAD_REQUEST", {
                  message: "User must agree to the TOS before signing up.",
                });
              }
              return {
                data: user,
              };
            },
          },
        },
      },
    });

#### [Using the Context Object](#using-the-context-object)

The context object (`ctx`), passed as the second argument to the hook, contains useful information. For `update` hooks, this includes the current `session`, which you can use to access the logged-in user's details.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      databaseHooks: {
        user: {
          update: {
            before: async (data, ctx) => {
              // You can access the session from the context object.
              if (ctx.context.session) {
                console.log("User update initiated by:", ctx.context.session.userId);
              }
              return { data };
            },
          },
        },
      },
    });

Much like standard hooks, database hooks also provide a `ctx` object that offers a variety of useful properties. Learn more in the [Hooks Documentation](https://www.better-auth.com/docs/concepts/hooks#ctx).

Plugins can define their own tables in the database to store additional data. They can also add columns to the core tables to store additional data. For example, the two factor authentication plugin adds the following columns to the `user` table:

*   `twoFactorEnabled`: Whether two factor authentication is enabled for the user.
*   `twoFactorSecret`: The secret key used to generate TOTP codes.
*   `twoFactorBackupCodes`: Encrypted backup codes for account recovery.

To add new tables and columns to your database, you have two options:

`CLI`: Use the migrate or generate command. These commands will scan your database and guide you through adding any missing tables or columns. `Manual Method`: Follow the instructions in the plugin documentation to manually add tables and columns.

Both methods ensure your database schema stays up to date with your plugins' requirements.

Since Better-Auth version `1.4` we've introduced experimental database joins support. This allows Better-Auth to perform multiple database queries in a single request, reducing the number of database roundtrips. Over 50 endpoints support joins, and we're constantly adding more.

Under the hood, our adapter system supports joins natively, meaning even if you don't enable experimental joins, it will still fallback to making multiple database queries and combining the results.

To enable joins, update your auth config with the following:

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

The Better-Auth `1.4` CLI will generate DrizzleORM and PrismaORM relationships for you so if you do not have those already be sure to update your schema by running our migrate or generate CLI commands to be up-to-date with the latest required schema.

It's very important to read the documentation regarding experimental joins for your given adapter:

*   [DrizzleORM](https://www.better-auth.com/docs/adapters/drizzle#joins-experimental)
*   [PrismaORM](https://www.better-auth.com/docs/adapters/prisma#joins-experimental)
*   [SQLite](https://www.better-auth.com/docs/adapters/sqlite#joins-experimental)
*   [MySQL](https://www.better-auth.com/docs/adapters/mysql#joins-experimental)
*   [PostgreSQL](https://www.better-auth.com/docs/adapters/postgresql#joins-experimental)
*   [MSSQL](https://www.better-auth.com/docs/adapters/mssql#joins-experimental)
*   [MongoDB](https://www.better-auth.com/docs/adapters/mongodb#joins-experimental)</content>
</page>

<page>
  <title>Passkey | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/passkey</url>
  <content>Passkeys are a secure, passwordless authentication method using cryptographic key pairs, supported by WebAuthn and FIDO2 standards in web browsers. They replace passwords with unique key pairs: a private key stored on the user's device and a public key shared with the website. Users can log in using biometrics, PINs, or security keys, providing strong, phishing-resistant authentication without traditional passwords.

The passkey plugin implementation is powered by [SimpleWebAuthn](https://simplewebauthn.dev/) behind the scenes.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

To add the passkey plugin to your auth config, you need to import the plugin and pass it to the `plugins` option of the auth instance.

auth.ts

    import { betterAuth } from "better-auth"
    import { passkey } from "@better-auth/passkey"
    
    export const auth = betterAuth({
        plugins: [ 
            passkey(), 
        ], 
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { passkeyClient } from "@better-auth/passkey/client"
    
    export const authClient = createAuthClient({
        plugins: [ 
            passkeyClient() 
        ] 
    })

### [Add/Register a passkey](#addregister-a-passkey)

To add or register a passkey make sure a user is authenticated and then call the `passkey.addPasskey` function provided by the client.

    const { data, error } = await authClient.passkey.addPasskey({    name: "example-passkey-name",    authenticatorAttachment: "cross-platform",});

| Prop | Description | Type |
| --- | --- | --- |
| `name?` | 
An optional name to label the authenticator account being registered. If not provided, it will default to the user's email address or user ID

 | `string` |
| `authenticatorAttachment?` | 

You can also specify the type of authenticator you want to register. Default behavior allows both platform and cross-platform passkeys

 | `"platform" | "cross-platform"` |

Setting `throw: true` in the fetch options has no effect for the register and sign-in passkey responses — they will always return a data object containing the error object.

### [Sign in with a passkey](#sign-in-with-a-passkey)

To sign in with a passkey you can use the `signIn.passkey` method. This will prompt the user to sign in with their passkey.

    const { data, error } = await authClient.signIn.passkey({    autoFill: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `autoFill?` | 
Browser autofill, a.k.a. Conditional UI. Read more: https://simplewebauthn.dev/docs/packages/browser#browser-autofill-aka-conditional-ui

 | `boolean` |

#### [Example Usage](#example-usage)

    // With post authentication redirect
    await authClient.signIn.passkey({
        autoFill: true,
        fetchOptions: {
            onSuccess(context) {
                // Redirect to dashboard after successful authentication
                window.location.href = "/dashboard";
            },
            onError(context) {
                // Handle authentication errors
                console.error("Authentication failed:", context.error.message);
            }
        }
    });

### [List passkeys](#list-passkeys)

You can list all of the passkeys for the authenticated user by calling `passkey.listUserPasskeys`:

GET

/passkey/list-user-passkeys

    const { data: passkeys, error } = await authClient.passkey.listUserPasskeys();

### [Deleting passkeys](#deleting-passkeys)

You can delete a passkey by calling `passkey.delete` and providing the passkey ID.

POST

/passkey/delete-passkey

    const { data, error } = await authClient.passkey.deletePasskey({    id: "some-passkey-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `id` | 
The ID of the passkey to delete.

 | `string` |

### [Updating passkey names](#updating-passkey-names)

POST

/passkey/update-passkey

    const { data, error } = await authClient.passkey.updatePasskey({    id: "id of passkey", // required    name: "my-new-passkey-name", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `id` | 
The ID of the passkey which you want to update.

 | `string` |
| `name` | 

The new name which the passkey will be updated to.

 | `string` |

### [Conditional UI](#conditional-ui)

The plugin supports conditional UI, which allows the browser to autofill the passkey if the user has already registered a passkey.

There are two requirements for conditional UI to work:

#### [Update input fields](#update-input-fields)

Add the `autocomplete` attribute with the value `webauthn` to your input fields. You can add this attribute to multiple input fields, but at least one is required for conditional UI to work.

The `webauthn` value should also be the last entry of the `autocomplete` attribute.

    <label for="name">Username:</label>
    <input type="text" name="name" autocomplete="username webauthn">
    <label for="password">Password:</label>
    <input type="password" name="password" autocomplete="current-password webauthn">

#### [Preload the passkeys](#preload-the-passkeys)

When your component mounts, you can preload the user's passkeys by calling the `authClient.signIn.passkey` method with the `autoFill` option set to `true`.

To prevent unnecessary calls, we will also add a check to see if the browser supports conditional UI.

Depending on the browser, a prompt will appear to autofill the passkey. If the user has multiple passkeys, they can select the one they want to use.

Some browsers also require the user to first interact with the input field before the autofill prompt appears.

### [Debugging](#debugging)

To test your passkey implementation you can use [emulated authenticators](https://developer.chrome.com/docs/devtools/webauthn). This way you can test the registration and sign-in process without even owning a physical device.

The plugin require a new table in the database to store passkey data.

Table Name: `passkey`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each passkey |
| name | string |  | The name of the passkey |
| publicKey | string | \- | The public key of the passkey |
| userId | string |  | The ID of the user |
| credentialID | string | \- | The unique identifier of the registered credential |
| counter | number | \- | The counter of the passkey |
| deviceType | string | \- | The type of device used to register the passkey |
| backedUp | boolean | \- | Whether the passkey is backed up |
| transports | string |  | The transports used to register the passkey |
| createdAt | Date |  | The time when the passkey was created |
| aaguid | string |  | Authenticator's Attestation GUID indicating the type of the authenticator |

**rpID**: A unique identifier for your website based on your auth server origin. `'localhost'` is okay for local dev. RP ID can be formed by discarding zero or more labels from the left of its effective domain until it hits an effective TLD. So `www.example.com` can use the RP IDs `www.example.com` or `example.com`. But not `com`, because that's an eTLD.

**rpName**: Human-readable title for your website.

**origin**: The origin URL at which your better-auth server is hosted. `http://localhost` and `http://localhost:PORT` are also valid. Do NOT include any trailing /.

**authenticatorSelection**: Allows customization of WebAuthn authenticator selection criteria. Leave unspecified for default settings.

*   `authenticatorAttachment`: Specifies the type of authenticator
    *   `platform`: Authenticator is attached to the platform (e.g., fingerprint reader)
    *   `cross-platform`: Authenticator is not attached to the platform (e.g., security key)
    *   Default: `not set` (both platform and cross-platform allowed, with platform preferred)
*   `residentKey`: Determines credential storage behavior.
    *   `required`: User MUST store credentials on the authenticator (highest security)
    *   `preferred`: Encourages credential storage but not mandatory
    *   `discouraged`: No credential storage required (fastest experience)
    *   Default: `preferred`
*   `userVerification`: Controls biometric/PIN verification during authentication:
    *   `required`: User MUST verify identity (highest security)
    *   `preferred`: Verification encouraged but not mandatory
    *   `discouraged`: No verification required (fastest experience)
    *   Default: `preferred`

**advanced**: Advanced options

*   `webAuthnChallengeCookie`: Cookie name for storing WebAuthn challenge ID during authentication flow (Default: `better-auth-passkey`)

When using the passkey plugin with Expo, you need to configure the `cookiePrefix` option in the Expo client to ensure passkey cookies are properly detected and stored.

By default, the passkey plugin uses `"better-auth-passkey"` as the challenge cookie name. Since this starts with `"better-auth"`, it will work with the default Expo client configuration. However, if you customize the `webAuthnChallengeCookie` option, you must also update the `cookiePrefix` in your Expo client configuration.

### [Example Configuration](#example-configuration)

If you're using a custom cookie name:

Server: auth.ts

    import { betterAuth } from "better-auth";
    import { passkey } from "@better-auth/passkey";
    
    export const auth = betterAuth({
        plugins: [
            passkey({
                advanced: {
                    webAuthnChallengeCookie: "my-app-passkey" // Custom cookie name
                }
            })
        ]
    });

Make sure to configure your Expo client with the matching prefix:

Client: auth-client.ts

    import { createAuthClient } from "better-auth/react";
    import { expoClient } from "@better-auth/expo/client";
    import { passkeyClient } from "@better-auth/passkey/client";
    import * as SecureStore from "expo-secure-store";
    
    export const authClient = createAuthClient({
        baseURL: "http://localhost:8081",
        plugins: [
            expoClient({
                storage: SecureStore,
                cookiePrefix: "my-app" // Must match the prefix of your custom cookie name
            }),
            passkeyClient()
        ]
    });

If you're using multiple authentication systems or custom cookie names, you can provide an array of prefixes:

Client: auth-client.ts

    expoClient({
        storage: SecureStore,
        cookiePrefix: ["better-auth", "my-app", "custom-auth"]
    })

If the `cookiePrefix` doesn't match the prefix of your `webAuthnChallengeCookie`, the passkey authentication flow will fail because the challenge cookie won't be stored and sent back to the server during verification.

For more information on Expo integration, see the [Expo documentation](https://www.better-auth.com/docs/integrations/expo).</content>
</page>

<page>
  <title>Username | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/username</url>
  <content>The username plugin is a lightweight plugin that adds username support to the email and password authenticator. This allows users to sign in and sign up with their username instead of their email.

### [Add Plugin to the server](#add-plugin-to-the-server)

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    export const auth = betterAuth({
        emailAndPassword: { 
            enabled: true, 
        }, 
        plugins: [ 
            username() 
        ] 
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { usernameClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [ 
            usernameClient() 
        ] 
    })

### [Sign up](#sign-up)

To sign up a user with username, you can use the existing `signUp.email` function provided by the client. The `signUp` function should take a new `username` property in the object.

    const { data, error } = await authClient.signUp.email({    email: "email@domain.com", // required    name: "Test User", // required    password: "password1234", // required    username: "test",    displayUsername: "Test User123",});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
The email of the user.

 | `string` |
| `name` | 

The name of the user.

 | `string` |
| `password` | 

The password of the user.

 | `string` |
| `username?` | 

The username of the user.

 | `string` |
| `displayUsername?` | 

An optional display username of the user.

 | `string` |

If only `username` is provided, the `displayUsername` will be set to the pre normalized version of the `username`. You can see the [Username Normalization](#username-normalization) and [Display Username Normalization](#display-username-normalization) sections for more details.

### [Sign in](#sign-in)

To sign in a user with username, you can use the `signIn.username` function provided by the client.

    const { data, error } = await authClient.signIn.username({    username: "test", // required    password: "password1234", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `username` | 
The username of the user.

 | `string` |
| `password` | 

The password of the user.

 | `string` |

### [Update username](#update-username)

To update the username of a user, you can use the `updateUser` function provided by the client.

    const { data, error } = await authClient.updateUser({    username: "new-username",});

| Prop | Description | Type |
| --- | --- | --- |
| `username?` | 
The username to update.

 | `string` |

### [Check if username is available](#check-if-username-is-available)

To check if a username is available, you can use the `isUsernameAvailable` function provided by the client.

POST

/is-username-available

    const { data: response, error } = await authClient.isUsernameAvailable({    username: "new-username", // required});if(response?.available) {    console.log("Username is available");} else {    console.log("Username is not available");}

| Prop | Description | Type |
| --- | --- | --- |
| `username` | 
The username to check.

 | `string` |

### [Min Username Length](#min-username-length)

The minimum length of the username. Default is `3`.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        plugins: [
            username({
                minUsernameLength: 5
            })
        ]
    })

### [Max Username Length](#max-username-length)

The maximum length of the username. Default is `30`.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        plugins: [
            username({
                maxUsernameLength: 100
            })
        ]
    })

### [Username Validator](#username-validator)

A function that validates the username. The function should return false if the username is invalid. By default, the username should only contain alphanumeric characters, underscores, and dots.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        plugins: [
            username({
                usernameValidator: (username) => {
                    if (username === "admin") {
                        return false
                    }
                    return true
                }
            })
        ]
    })

### [Display Username Validator](#display-username-validator)

A function that validates the display username. The function should return false if the display username is invalid. By default, no validation is applied to display username.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        plugins: [
            username({
                displayUsernameValidator: (displayUsername) => {
                    // Allow only alphanumeric characters, underscores, and hyphens
                    return /^[a-zA-Z0-9_-]+$/.test(displayUsername)
                }
            })
        ]
    })

### [Username Normalization](#username-normalization)

A function that normalizes the username, or `false` if you want to disable normalization.

By default, usernames are normalized to lowercase, so "TestUser" and "testuser", for example, are considered the same username. The `username` field will contain the normalized (lower case) username, while `displayUsername` will contain the original `username`.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        plugins: [
            username({
                usernameNormalization: (username) => {
                    return username.toLowerCase()
                        .replaceAll("0", "o")
                        .replaceAll("3", "e")
                        .replaceAll("4", "a");
                }
            })
        ]
    })

### [Display Username Normalization](#display-username-normalization)

A function that normalizes the display username, or `false` to disable normalization.

By default, display usernames are not normalized. When only `username` is provided during signup or update, the `displayUsername` will be set to match the original `username` value (before normalization). You can also explicitly set a `displayUsername` which will be preserved as-is. For custom normalization, provide a function that takes the display username as input and returns the normalized version.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        plugins: [
            username({
                displayUsernameNormalization: (displayUsername) => displayUsername.toLowerCase(),
            })
        ]
    })

### [Validation Order](#validation-order)

By default, username and display username are validated before normalization. You can change this behavior by setting `validationOrder` to `post-normalization`.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        plugins: [
            username({
                validationOrder: {
                    username: "post-normalization",
                    displayUsername: "post-normalization",
                }
            })
        ]
    })

### [Disable Is Username Available](#disable-is-username-available)

By default, the plugin exposes an endpoint `/is-username-available` to check if a username is available. You can disable this endpoint by providing `disablePaths` option to the better-auth configuration. This is useful if you want to protect usernames from being enumerated.

auth.ts

    import { betterAuth } from "better-auth"
    import { username } from "better-auth/plugins"
    
    const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
        },
        disablePaths: ["/is-username-available"],
        plugins: [
            username()
        ]
    })

The plugin requires 2 fields to be added to the user table:

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| username | string | \- | The username of the user |
| displayUsername | string | \- | Non normalized username of the user |

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/username.mdx)</content>
</page>

<page>
  <title>Magic link | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/magic-link</url>
  <content>Magic link or email link is a way to authenticate users without a password. When a user enters their email, a link is sent to their email. When the user clicks on the link, they are authenticated.

### [Add the server Plugin](#add-the-server-plugin)

Add the magic link plugin to your server:

server.ts

    import { betterAuth } from "better-auth";
    import { magicLink } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [
            magicLink({
                sendMagicLink: async ({ email, token, url }, ctx) => {
                    // send email to user
                }
            })
        ]
    })

### [Add the client Plugin](#add-the-client-plugin)

Add the magic link plugin to your client:

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { magicLinkClient } from "better-auth/client/plugins";
    export const authClient = createAuthClient({
        plugins: [
            magicLinkClient()
        ]
    });

### [Sign In with Magic Link](#sign-in-with-magic-link)

To sign in with a magic link, you need to call `signIn.magicLink` with the user's email address. The `sendMagicLink` function is called to send the magic link to the user's email.

    const { data, error } = await authClient.signIn.magicLink({    email: "user@email.com", // required    name: "my-name",    callbackURL: "/dashboard",    newUserCallbackURL: "/welcome",    errorCallbackURL: "/error",});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to send the magic link.

 | `string` |
| `name?` | 

User display name. Only used if the user is registering for the first time.

 | `string` |
| `callbackURL?` | 

URL to redirect after magic link verification.

 | `string` |
| `newUserCallbackURL?` | 

URL to redirect after new user signup

 | `string` |
| `errorCallbackURL?` | 

URL to redirect if an error happen on verification If only callbackURL is provided but without an `errorCallbackURL` then they will be redirected to the callbackURL with an `error` query parameter.

 | `string` |

If the user has not signed up, unless `disableSignUp` is set to `true`, the user will be signed up automatically.

### [Verify Magic Link](#verify-magic-link)

When you send the URL generated by the `sendMagicLink` function to a user, clicking the link will authenticate them and redirect them to the `callbackURL` specified in the `signIn.magicLink` function. If an error occurs, the user will be redirected to the `callbackURL` with an error query parameter.

If no `callbackURL` is provided, the user will be redirected to the root URL.

If you want to handle the verification manually, (e.g, if you send the user a different URL), you can use the `verify` function.

    const { data, error } = await authClient.magicLink.verify({    query: {        token: "123456", // required        callbackURL: "/dashboard",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `token` | 
Verification token.

 | `string` |
| `callbackURL?` | 

URL to redirect after magic link verification, if not provided will return the session.

 | `string` |

**sendMagicLink**: The `sendMagicLink` function is called when a user requests a magic link. It takes an object with the following properties:

*   `email`: The email address of the user.
*   `url`: The URL to be sent to the user. This URL contains the token.
*   `token`: The token if you want to send the token with custom URL.

and a `ctx` context object as the second parameter.

**expiresIn**: specifies the time in seconds after which the magic link will expire. The default value is `300` seconds (5 minutes).

**disableSignUp**: If set to `true`, the user will not be able to sign up using the magic link. The default value is `false`.

**generateToken**: The `generateToken` function is called to generate a token which is used to uniquely identify the user. The default value is a random string. There is one parameter:

*   `email`: The email address of the user.

When using `generateToken`, ensure that the returned string is hard to guess because it is used to verify who someone actually is in a confidential way. By default, we return a long and cryptographically secure string.

**storeToken**: The `storeToken` function is called to store the magic link token in the database. The default value is `"plain"`.

The `storeToken` function can be one of the following:

*   `"plain"`: The token is stored in plain text.
*   `"hashed"`: The token is hashed using the default hasher.
*   `{ type: "custom-hasher", hash: (token: string) => Promise<string> }`: The token is hashed using a custom hasher.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/magic-link.mdx)</content>
</page>

<page>
  <title>Better Auth Fastify Integration Guide | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/fastify</url>
  <content>    import Fastify from "fastify";
    import { auth } from "./auth"; // Your configured Better Auth instance
    
    const fastify = Fastify({ logger: true });
    
    // Register authentication endpoint
    fastify.route({
      method: ["GET", "POST"],
      url: "/api/auth/*",
      async handler(request, reply) {
        try {
          // Construct request URL
          const url = new URL(request.url, `http://${request.headers.host}`);
          
          // Convert Fastify headers to standard Headers object
          const headers = new Headers();
          Object.entries(request.headers).forEach(([key, value]) => {
            if (value) headers.append(key, value.toString());
          });
    
          // Create Fetch API-compatible request
          const req = new Request(url.toString(), {
            method: request.method,
            headers,
            body: request.body ? JSON.stringify(request.body) : undefined,
          });
    
          // Process authentication request
          const response = await auth.handler(req);
    
          // Forward response to client
          reply.status(response.status);
          response.headers.forEach((value, key) => reply.header(key, value));
          reply.send(response.body ? await response.text() : null);
    
        } catch (error) {
          fastify.log.error("Authentication Error:", error);
          reply.status(500).send({ 
            error: "Internal authentication error",
            code: "AUTH_FAILURE"
          });
        }
      }
    });
    
    // Initialize server
    fastify.listen({ port: 4000 }, (err) => {
      if (err) {
        fastify.log.error(err);
        process.exit(1);
      }
      console.log("Server running on port 4000");
    });</content>
</page>

<page>
  <title>Email OTP | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/email-otp</url>
  <content>The Email OTP plugin allows user to sign in, verify their email, or reset their password using a one-time password (OTP) sent to their email address.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

Add the `emailOTP` plugin to your auth config and implement the `sendVerificationOTP()` method.

auth.ts

    import { betterAuth } from "better-auth"
    import { emailOTP } from "better-auth/plugins"
    
    export const auth = betterAuth({
        // ... other config options
        plugins: [
            emailOTP({ 
                async sendVerificationOTP({ email, otp, type }) { 
                    if (type === "sign-in") { 
                        // Send the OTP for sign in
                    } else if (type === "email-verification") { 
                        // Send the OTP for email verification
                    } else { 
                        // Send the OTP for password reset
                    } 
                }, 
            }) 
        ]
    })

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { emailOTPClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            emailOTPClient()
        ]
    })

### [Send an OTP](#send-an-otp)

Use the `sendVerificationOtp()` method to send an OTP to the user's email address.

POST

/email-otp/send-verification-otp

    const { data, error } = await authClient.emailOtp.sendVerificationOtp({    email: "user@example.com", // required    type: "sign-in", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to send the OTP.

 | `string` |
| `type` | 

Type of the OTP. `sign-in`, `email-verification`, or `forget-password`.

 | `"email-verification" | "sign-in" | "forget-password"` |

### [Check an OTP (optional)](#check-an-otp-optional)

Use the `checkVerificationOtp()` method to check if an OTP is valid.

POST

/email-otp/check-verification-otp

    const { data, error } = await authClient.emailOtp.checkVerificationOtp({    email: "user@example.com", // required    type: "sign-in", // required    otp: "123456", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to send the OTP.

 | `string` |
| `type` | 

Type of the OTP. `sign-in`, `email-verification`, or `forget-password`.

 | `"email-verification" | "sign-in" | "forget-password"` |
| `otp` | 

OTP sent to the email.

 | `string` |

### [Sign In with OTP](#sign-in-with-otp)

To sign in with OTP, use the `sendVerificationOtp()` method to send a "sign-in" OTP to the user's email address.

POST

/email-otp/send-verification-otp

    const { data, error } = await authClient.emailOtp.sendVerificationOtp({    email: "user@example.com", // required    type: "sign-in", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to send the OTP.

 | `string` |
| `type` | 

Type of the OTP.

 | `"sign-in"` |

Once the user provides the OTP, you can sign in the user using the `signIn.emailOtp()` method.

    const { data, error } = await authClient.signIn.emailOtp({    email: "user@example.com", // required    otp: "123456", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to sign in.

 | `string` |
| `otp` | 

OTP sent to the email.

 | `string` |

If the user is not registered, they'll be automatically registered. If you want to prevent this, you can pass `disableSignUp` as `true` in the [options](#options).

### [Verify Email with OTP](#verify-email-with-otp)

To verify the user's email address with OTP, use the `sendVerificationOtp()` method to send an "email-verification" OTP to the user's email address.

POST

/email-otp/send-verification-otp

    const { data, error } = await authClient.emailOtp.sendVerificationOtp({    email: "user@example.com", // required    type: "email-verification", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to send the OTP.

 | `string` |
| `type` | 

Type of the OTP.

 | `"email-verification"` |

Once the user provides the OTP, use the `verifyEmail()` method to complete email verification.

POST

/email-otp/verify-email

    const { data, error } = await authClient.emailOtp.verifyEmail({    email: "user@example.com", // required    otp: "123456", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to verify.

 | `string` |
| `otp` | 

OTP to verify.

 | `string` |

### [Reset Password with OTP](#reset-password-with-otp)

To reset the user's password with OTP, use the `forgetPassword.emailOTP()` method to send a "forget-password" OTP to the user's email address.

POST

/forget-password/email-otp

    const { data, error } = await authClient.forgetPassword.emailOtp({    email: "user@example.com", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to send the OTP.

 | `string` |

Once the user provides the OTP, use the `checkVerificationOtp()` method to check if it's valid (optional).

POST

/email-otp/check-verification-otp

    const { data, error } = await authClient.emailOtp.checkVerificationOtp({    email: "user@example.com", // required    type: "forget-password", // required    otp: "123456", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to send the OTP.

 | `string` |
| `type` | 

Type of the OTP.

 | `"forget-password"` |
| `otp` | 

OTP sent to the email.

 | `string` |

Then, use the `resetPassword()` method to reset the user's password.

POST

/email-otp/reset-password

    const { data, error } = await authClient.emailOtp.resetPassword({    email: "user@example.com", // required    otp: "123456", // required    password: "new-secure-password", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
Email address to reset the password.

 | `string` |
| `otp` | 

OTP sent to the email.

 | `string` |
| `password` | 

New password.

 | `string` |

### [Override Default Email Verification](#override-default-email-verification)

To override the default email verification, pass `overrideDefaultEmailVerification: true` in the options. This will make the system use an email OTP instead of the default verification link whenever email verification is triggered. In other words, the user will verify their email using an OTP rather than clicking a link.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      plugins: [
        emailOTP({
          overrideDefaultEmailVerification: true, 
          async sendVerificationOTP({ email, otp, type }) {
            // Implement the sendVerificationOTP method to send the OTP to the user's email address
          },
        }),
      ],
    });

*   `sendVerificationOTP`: A function that sends the OTP to the user's email address. The function receives an object with the following properties:
    
    *   `email`: The user's email address.
    *   `otp`: The OTP to send.
    *   `type`: The type of OTP to send. Can be "sign-in", "email-verification", or "forget-password".
    
    It is recommended to not await the email sending to avoid timing attacks. On serverless platforms, use `waitUntil` or similar to ensure the email is sent.
    
*   `otpLength`: The length of the OTP. Defaults to `6`.
    
*   `expiresIn`: The expiry time of the OTP in seconds. Defaults to `300` seconds.
    

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        plugins: [
            emailOTP({
                otpLength: 8,
                expiresIn: 600
            })
        ]
    })

*   `sendVerificationOnSignUp`: A boolean value that determines whether to send the OTP when a user signs up. Defaults to `false`.
    
*   `disableSignUp`: A boolean value that determines whether to prevent automatic sign-up when the user is not registered. Defaults to `false`.
    
*   `generateOTP`: A function that generates the OTP. Defaults to a random 6-digit number.
    
*   `allowedAttempts`: The maximum number of attempts allowed for verifying an OTP. Defaults to `3`. After exceeding this limit, the OTP becomes invalid and the user needs to request a new one.
    

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        plugins: [
            emailOTP({
                allowedAttempts: 5, // Allow 5 attempts before invalidating the OTP
                expiresIn: 300
            })
        ]
    })

When the maximum attempts are exceeded, the `verifyOTP`, `signIn.emailOtp`, `verifyEmail`, and `resetPassword` methods will return an error with code `TOO_MANY_ATTEMPTS`.

*   `storeOTP`: The method to store the OTP in your database, whether `encrypted`, `hashed` or `plain` text. Default is `plain` text.

Note: This will not affect the OTP sent to the user, it will only affect the OTP stored in your database.

Alternatively, you can pass a custom encryptor or hasher to store the OTP in your database.

**Custom encryptor**

auth.ts

    emailOTP({
        storeOTP: { 
            encrypt: async (otp) => {
                return myCustomEncryptor(otp);
            },
            decrypt: async (otp) => {
                return myCustomDecryptor(otp);
            },
        }
    })

**Custom hasher**

auth.ts

    emailOTP({
        storeOTP: {
            hash: async (otp) => {
                return myCustomHasher(otp);
            },
        }
    })

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/email-otp.mdx)</content>
</page>

<page>
  <title>Next.js integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/next#server-action-cookies</url>
  <content>Better Auth can be easily integrated with Next.js. Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Create API Route](#create-api-route)

We need to mount the handler to an API route. Create a route file inside `/api/auth/[...all]` directory. And add the following code:

api/auth/\[...all\]/route.ts

    import { auth } from "@/lib/auth";
    import { toNextJsHandler } from "better-auth/next-js";
    
    export const { GET, POST } = toNextJsHandler(auth);

You can change the path on your better-auth configuration but it's recommended to keep it as `/api/auth/[...all]`

For `pages` route, you need to use `toNodeHandler` instead of `toNextJsHandler` and set `bodyParser` to `false` in the `config` object. Here is an example:

pages/api/auth/\[...all\].ts

    import { toNodeHandler } from "better-auth/node"
    import { auth } from "@/lib/auth"
    
    // Disallow body parsing, we will parse it manually
    export const config = { api: { bodyParser: false } }
    
    export default toNodeHandler(auth.handler)

Create a client instance. You can name the file anything you want. Here we are creating `client.ts` file inside the `lib/` directory.

auth-client.ts

    import { createAuthClient } from "better-auth/react" // make sure to import from better-auth/react
    
    export const authClient =  createAuthClient({
        //you can pass client configuration here
    })

Once you have created the client, you can use it to sign up, sign in, and perform other actions. Some of the actions are reactive. The client uses [nano-store](https://github.com/nanostores/nanostores) to store the state and re-render the components when the state changes.

The client also uses [better-fetch](https://github.com/bekacru/better-fetch) to make the requests. You can pass the fetch configuration to the client.

The `api` object exported from the auth instance contains all the actions that you can perform on the server. Every endpoint made inside Better Auth is a invocable as a function. Including plugins endpoints.

**Example: Getting Session on a server action**

server.ts

    import { auth } from "@/lib/auth"
    import { headers } from "next/headers"
    
    const someAuthenticatedAction = async () => {
        "use server";
        const session = await auth.api.getSession({
            headers: await headers()
        })
    };

**Example: Getting Session on a RSC**

    import { auth } from "@/lib/auth"
    import { headers } from "next/headers"
    
    export async function ServerComponent() {
        const session = await auth.api.getSession({
            headers: await headers()
        })
        if(!session) {
            return <div>Not authenticated</div>
        }
        return (
            <div>
                <h1>Welcome {session.user.name}</h1>
            </div>
        )
    }

As RSCs cannot set cookies, the [cookie cache](https://www.better-auth.com/docs/concepts/session-management#cookie-cache) will not be refreshed until the server is interacted with from the client via Server Actions or Route Handlers.

### [Server Action Cookies](#server-action-cookies)

When you call a function that needs to set cookies, like `signInEmail` or `signUpEmail` in a server action, cookies won’t be set. This is because server actions need to use the `cookies` helper from Next.js to set cookies.

To simplify this, you can use the `nextCookies` plugin, which will automatically set cookies for you whenever a `Set-Cookie` header is present in the response.

auth.ts

    import { betterAuth } from "better-auth";
    import { nextCookies } from "better-auth/next-js";
    
    export const auth = betterAuth({
        //...your config
        plugins: [nextCookies()] // make sure this is the last plugin in the array
    })

Now, when you call functions that set cookies, they will be automatically set.

    "use server";
    import { auth } from "@/lib/auth"
    
    const signIn = async () => {
        await auth.api.signInEmail({
            body: {
                email: "user@email.com",
                password: "password",
            }
        })
    }

In Next.js proxy/middleware, it's recommended to only check for the existence of a session cookie to handle redirection. To avoid blocking requests by making API or database calls.

### [Next.js 16+ (Proxy)](#nextjs-16-proxy)

Next.js 16 replaces "middleware" with "proxy". You can use the Node.js runtime for full session validation with database checks:

proxy.ts

    import { NextRequest, NextResponse } from "next/server";
    import { headers } from "next/headers";
    import { auth } from "@/lib/auth";
    
    export async function proxy(request: NextRequest) {
        const session = await auth.api.getSession({
            headers: await headers()
        })
    
        // THIS IS NOT SECURE!
        // This is the recommended approach to optimistically redirect users
        // We recommend handling auth checks in each page/route
        if(!session) {
            return NextResponse.redirect(new URL("/sign-in", request.url));
        }
    
        return NextResponse.next();
    }
    
    export const config = {
      matcher: ["/dashboard"], // Specify the routes the middleware applies to
    };

For cookie-only checks (faster but less secure), use `getSessionCookie`:

proxy.ts

    import { NextRequest, NextResponse } from "next/server";
    import { getSessionCookie } from "better-auth/cookies";
    
    export async function proxy(request: NextRequest) {
    	const sessionCookie = getSessionCookie(request);
    
        // THIS IS NOT SECURE!
        // This is the recommended approach to optimistically redirect users
        // We recommend handling auth checks in each page/route
    	if (!sessionCookie) {
    		return NextResponse.redirect(new URL("/", request.url));
    	}
    
    	return NextResponse.next();
    }
    
    export const config = {
    	matcher: ["/dashboard"], // Specify the routes the middleware applies to
    };

**Migration from middleware:** Rename `middleware.ts` → `proxy.ts` and `middleware` → `proxy` function. All Better Auth methods work identically.

### [Next.js 15.2.0+ (Node.js Runtime Middleware)](#nextjs-1520-nodejs-runtime-middleware)

From Next.js 15.2.0, you can use the Node.js runtime in middleware for full session validation with database checks:

middleware.ts

    import { NextRequest, NextResponse } from "next/server";
    import { headers } from "next/headers";
    import { auth } from "@/lib/auth";
    
    export async function middleware(request: NextRequest) {
        const session = await auth.api.getSession({
            headers: await headers()
        })
    
        // THIS IS NOT SECURE!
        // This is the recommended approach to optimistically redirect users
        // We recommend handling auth checks in each page/route
        if(!session) {
            return NextResponse.redirect(new URL("/sign-in", request.url));
        }
    
        return NextResponse.next();
    }
    
    export const config = {
      runtime: "nodejs", // Required for auth.api calls
      matcher: ["/dashboard"], // Specify the routes the middleware applies to
    };

Node.js runtime in middleware is experimental in Next.js versions before 16. Consider upgrading to Next.js 16+ for stable proxy support.

### [Next.js 13-15.1.x (Edge Runtime Middleware)](#nextjs-13-151x-edge-runtime-middleware)

In older Next.js versions, middleware runs on the Edge Runtime and cannot make database calls. Use cookie-based checks for optimistic redirects:

The `getSessionCookie()` function does not automatically reference the auth config specified in `auth.ts`. Therefore, if you customized the cookie name or prefix, you need to ensure that the configuration in `getSessionCookie()` matches the config defined in your `auth.ts`.

#### [For Next.js release `15.1.7` and below](#for-nextjs-release-1517-and-below)

If you need the full session object, you'll have to fetch it from the `/api/auth/get-session` API route. Since Next.js middleware doesn't support running Node.js APIs directly, you must make an HTTP request.

The example uses [better-fetch](https://better-fetch.vercel.app/), but you can use any fetch library.

middleware.ts

    import { betterFetch } from "@better-fetch/fetch";
    import type { auth } from "@/lib/auth";
    import { NextRequest, NextResponse } from "next/server";
    
    type Session = typeof auth.$Infer.Session;
    
    export async function middleware(request: NextRequest) {
    	const { data: session } = await betterFetch<Session>("/api/auth/get-session", {
    		baseURL: request.nextUrl.origin,
    		headers: {
    			cookie: request.headers.get("cookie") || "", // Forward the cookies from the request
    		},
    	});
    
    	if (!session) {
    		return NextResponse.redirect(new URL("/sign-in", request.url));
    	}
    
    	return NextResponse.next();
    }
    
    export const config = {
    	matcher: ["/dashboard"], // Apply middleware to specific routes
    };

#### [For Next.js release `15.2.0` and above](#for-nextjs-release-1520-and-above)

From Next.js 15.2.0, you can use the Node.js runtime in middleware for full session validation with database checks:

You may refer to the [Next.js documentation](https://nextjs.org/docs/app/building-your-application/routing/middleware#runtime) for more information about runtime configuration, and how to enable it. Be careful when using the new runtime. It's an experimental feature and it may be subject to breaking changes.

middleware.ts

    import { NextRequest, NextResponse } from "next/server";
    import { headers } from "next/headers";
    import { auth } from "@/lib/auth";
    
    export async function middleware(request: NextRequest) {
        const session = await auth.api.getSession({
            headers: await headers()
        })
    
        if(!session) {
            return NextResponse.redirect(new URL("/sign-in", request.url));
        }
    
        return NextResponse.next();
    }
    
    export const config = {
      runtime: "nodejs",
      matcher: ["/dashboard"], // Apply middleware to specific routes
    };

#### [Cookie-based checks (recommended for all versions)](#cookie-based-checks-recommended-for-all-versions)

middleware.ts

    import { NextRequest, NextResponse } from "next/server";
    import { getSessionCookie } from "better-auth/cookies";
    
    export async function middleware(request: NextRequest) {
    	const sessionCookie = getSessionCookie(request);
    
        // THIS IS NOT SECURE!
        // This is the recommended approach to optimistically redirect users
        // We recommend handling auth checks in each page/route
    	if (!sessionCookie) {
    		return NextResponse.redirect(new URL("/", request.url));
    	}
    
    	return NextResponse.next();
    }
    
    export const config = {
    	matcher: ["/dashboard"], // Specify the routes the middleware applies to
    };

**Security Warning:** The `getSessionCookie` function only checks for the existence of a session cookie; it does **not** validate it. Relying solely on this check for security is dangerous, as anyone can manually create a cookie to bypass it. You must always validate the session on your server for any protected actions or pages.

If you have a custom cookie name or prefix, you can pass it to the `getSessionCookie` function.

    const sessionCookie = getSessionCookie(request, {
        cookieName: "my_session_cookie",
        cookiePrefix: "my_prefix"
    });

Alternatively, you can use the `getCookieCache` helper to get the session object from the cookie cache.

middleware.ts

    import { getCookieCache } from "better-auth/cookies";
    
    export async function middleware(request: NextRequest) {
    	const session = await getCookieCache(request);
    	if (!session) {
    		return NextResponse.redirect(new URL("/sign-in", request.url));
    	}
    	return NextResponse.next();
    }

### [How to handle auth checks in each page/route](#how-to-handle-auth-checks-in-each-pageroute)

In this example, we are using the `auth.api.getSession` function within a server component to get the session object, then we are checking if the session is valid. If it's not, we are redirecting the user to the sign-in page.

app/dashboard/page.tsx

    import { auth } from "@/lib/auth";
    import { headers } from "next/headers";
    import { redirect } from "next/navigation";
    
    export default async function DashboardPage() {
        const session = await auth.api.getSession({
            headers: await headers()
        })
    
        if(!session) {
            redirect("/sign-in")
        }
    
        return <h1>Welcome {session.user.name}</h1>
    }

Better Auth is fully compatible with Next.js 16. The main change is that "middleware" is now called "proxy". See the [Auth Protection](#auth-protection) section above for Next.js 16+ proxy examples.

### [Migration Guide](#migration-guide)

Use Next.js codemod for automatic migration:

    npx @next/codemod@canary middleware-to-proxy .

Or manually:

*   Rename `middleware.ts` → `proxy.ts`
*   Change function name: `middleware` → `proxy`

All Better Auth methods work identically. See the [Next.js migration guide](https://nextjs.org/docs/app/api-reference/file-conventions/proxy#migration-to-proxy) for details.</content>
</page>

<page>
  <title>Two-Factor Authentication (2FA) | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/2fa#schema</url>
  <content>`OTP` `TOTP` `Backup Codes` `Trusted Devices`

Two-Factor Authentication (2FA) adds an extra security step when users log in. Instead of just using a password, they'll need to provide a second form of verification. This makes it much harder for unauthorized people to access accounts, even if they've somehow gotten the password.

This plugin offers two main methods to do a second factor verification:

1.  **OTP (One-Time Password)**: A temporary code sent to the user's email or phone.
2.  **TOTP (Time-based One-Time Password)**: A code generated by an app on the user's device.

**Additional features include:**

*   Generating backup codes for account recovery
*   Enabling/disabling 2FA
*   Managing trusted devices

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

Add the two-factor plugin to your auth configuration and specify your app name as the issuer.

auth.ts

    import { betterAuth } from "better-auth"
    import { twoFactor } from "better-auth/plugins"
    
    export const auth = betterAuth({
        // ... other config options
        appName: "My App", // provide your app name. It'll be used as an issuer.
        plugins: [
            twoFactor() 
        ]
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

Add the client plugin and Specify where the user should be redirected if they need to verify 2nd factor

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { twoFactorClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            twoFactorClient()
        ]
    })

### [Enabling 2FA](#enabling-2fa)

To enable two-factor authentication, call `twoFactor.enable` with the user's password and issuer (optional):

    const { data, error } = await authClient.twoFactor.enable({    password: "secure-password", // required    issuer: "my-app-name",});

| Prop | Description | Type |
| --- | --- | --- |
| `password` | 
The user's password

 | `string` |
| `issuer?` | 

An optional custom issuer for the TOTP URI. Defaults to app-name defined in your auth config.

 | `string` |

When 2FA is enabled:

*   An encrypted `secret` and `backupCodes` are generated.
*   `enable` returns `totpURI` and `backupCodes`.

Note: `twoFactorEnabled` won’t be set to `true` until the user verifies their TOTP code. Learn more about verifying TOTP [here](#totp). You can skip verification by setting `skipVerificationOnEnable` to true in your plugin config.

Two Factor can only be enabled for credential accounts at the moment. For social accounts, it's assumed the provider already handles 2FA.

### [Sign In with 2FA](#sign-in-with-2fa)

When a user with 2FA enabled tries to sign in via email, the response object will contain `twoFactorRedirect` set to `true`. This indicates that the user needs to verify their 2FA code.

You can handle this in the `onSuccess` callback or by providing a `onTwoFactorRedirect` callback in the plugin config.

sign-in.tsx

    await authClient.signIn.email({
            email: "user@example.com",
            password: "password123",
        },
        {
            async onSuccess(context) {
                if (context.data.twoFactorRedirect) {
                    // Handle the 2FA verification in place
                }
            },
        }
    )

Using the `onTwoFactorRedirect` config:

sign-in.ts

    import { createAuthClient } from "better-auth/client";
    import { twoFactorClient } from "better-auth/client/plugins";
    
    const authClient = createAuthClient({
        plugins: [
            twoFactorClient({
                onTwoFactorRedirect(){
                    // Handle the 2FA verification globally
                },
            }),
        ],
    });

**With `auth.api`**

When you call `auth.api.signInEmail` on the server, and the user has 2FA enabled, it will return an object where `twoFactorRedirect` is set to `true`. This behavior isn’t inferred in TypeScript, which can be misleading. You can check using `in` instead to check if `twoFactorRedirect` is set to `true`.

    const response = await auth.api.signInEmail({
    	body: {
    		email: "test@test.com",
    		password: "test",
    	},
    });
    
    if ("twoFactorRedirect" in response) {
    	// Handle the 2FA verification in place
    }

### [Disabling 2FA](#disabling-2fa)

To disable two-factor authentication, call `twoFactor.disable` with the user's password:

    const { data, error } = await authClient.twoFactor.disable({    password, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `password` | 
The user's password

 | `string` |

### [TOTP](#totp)

TOTP (Time-Based One-Time Password) is an algorithm that generates a unique password for each login attempt using time as a counter. Every fixed interval (Better Auth defaults to 30 seconds), a new password is generated. This addresses several issues with traditional passwords: they can be forgotten, stolen, or guessed. OTPs solve some of these problems, but their delivery via SMS or email can be unreliable (or even risky, considering it opens new attack vectors).

TOTP, however, generates codes offline, making it both secure and convenient. You just need an authenticator app on your phone.

#### [Getting TOTP URI](#getting-totp-uri)

After enabling 2FA, you can get the TOTP URI to display to the user. This URI is generated by the server using the `secret` and `issuer` and can be used to generate a QR code for the user to scan with their authenticator app.

POST

/two-factor/get-totp-uri

    const { data, error } = await authClient.twoFactor.getTotpUri({    password, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `password` | 
The user's password

 | `string` |

**Example: Using React**

Once you have the TOTP URI, you can use it to generate a QR code for the user to scan with their authenticator app.

user-card.tsx

    import QRCode from "react-qr-code";
    
    export default function UserCard({ password }: { password: string }){
        const { data: session } = client.useSession();
    	const { data: qr } = useQuery({
    		queryKey: ["two-factor-qr"],
    		queryFn: async () => {
    			const res = await authClient.twoFactor.getTotpUri({ password });
    			return res.data;
    		},
    		enabled: !!session?.user.twoFactorEnabled,
    	});
        return (
            <QRCode value={qr?.totpURI || ""} />
       )
    }

By default the issuer for TOTP is set to the app name provided in the auth config or if not provided it will be set to `Better Auth`. You can override this by passing `issuer` to the plugin config.

#### [Verifying TOTP](#verifying-totp)

After the user has entered their 2FA code, you can verify it using `twoFactor.verifyTotp` method. `Better Auth` follows standard practice by accepting TOTP codes from one period before and one after the current code, ensuring users can authenticate even with minor time delays on their end.

POST

/two-factor/verify-totp

    const { data, error } = await authClient.twoFactor.verifyTotp({    code: "012345", // required    trustDevice: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `code` | 
The otp code to verify.

 | `string` |
| `trustDevice?` | 

If true, the device will be trusted for 30 days. It'll be refreshed on every sign in request within this time.

 | `boolean` |

### [OTP](#otp)

OTP (One-Time Password) is similar to TOTP but a random code is generated and sent to the user's email or phone.

Before using OTP to verify the second factor, you need to configure `sendOTP` in your Better Auth instance. This function is responsible for sending the OTP to the user's email, phone, or any other method supported by your application.

auth.ts

    import { betterAuth } from "better-auth"
    import { twoFactor } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [
            twoFactor({
              	otpOptions: {
    				async sendOTP({ user, otp }, ctx) {
                        // send otp to user
    				},
    			},
            })
        ]
    })

#### [Sending OTP](#sending-otp)

Sending an OTP is done by calling the `twoFactor.sendOtp` function. This function will trigger your sendOTP implementation that you provided in the Better Auth configuration.

    const { data, error } = await authClient.twoFactor.sendOtp({    trustDevice: true,});if (data) {    // redirect or show the user to enter the code}

| Prop | Description | Type |
| --- | --- | --- |
| `trustDevice?` | 
If true, the device will be trusted for 30 days. It'll be refreshed on every sign in request within this time.

 | `boolean` |

#### [Verifying OTP](#verifying-otp)

After the user has entered their OTP code, you can verify it

POST

/two-factor/verify-otp

    const { data, error } = await authClient.twoFactor.verifyOtp({    code: "012345", // required    trustDevice: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `code` | 
The otp code to verify.

 | `string` |
| `trustDevice?` | 

If true, the device will be trusted for 30 days. It'll be refreshed on every sign in request within this time.

 | `boolean` |

### [Backup Codes](#backup-codes)

Backup codes are generated and stored in the database. This can be used to recover access to the account if the user loses access to their phone or email.

#### [Generating Backup Codes](#generating-backup-codes)

Generate backup codes for account recovery:

POST

/two-factor/generate-backup-codes

    const { data, error } = await authClient.twoFactor.generateBackupCodes({    password, // required});if (data) {    // Show the backup codes to the user}

| Prop | Description | Type |
| --- | --- | --- |
| `password` | 
The users password.

 | `string` |

When you generate backup codes, the old backup codes will be deleted and new ones will be generated.

#### [Using Backup Codes](#using-backup-codes)

You can now allow users to provider backup code as account recover method.

POST

/two-factor/verify-backup-code

    const { data, error } = await authClient.twoFactor.verifyBackupCode({    code: "123456", // required    disableSession: false,    trustDevice: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `code` | 
A backup code to verify.

 | `string` |
| `disableSession?` | 

If true, the session cookie will not be set.

 | `boolean` |
| `trustDevice?` | 

If true, the device will be trusted for 30 days. It'll be refreshed on every sign in request within this time.

 | `boolean` |

Once a backup code is used, it will be removed from the database and can't be used again.

#### [Viewing Backup Codes](#viewing-backup-codes)

To display the backup codes to the user, you can call `viewBackupCodes` on the server. This will return the backup codes in the response. You should only this if the user has a fresh session - a session that was just created.

POST

/two-factor/view-backup-codes

    const data = await auth.api.viewBackupCodes({    body: {        userId: "user-id",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `userId?` | 
The user ID to view all backup codes.

 | `string | null` |

### [Trusted Devices](#trusted-devices)

You can mark a device as trusted by passing `trustDevice` to `verifyTotp` or `verifyOtp`.

    const verify2FA = async (code: string) => {
        const { data, error } = await authClient.twoFactor.verifyTotp({
            code,
            trustDevice: true, // Mark this device as trusted
        })
        if (data) {
            // 2FA verified and device trusted
        }
    }

When `trustDevice` is set to `true`, the current device will be remembered for 30 days. During this period, the user won't be prompted for 2FA on subsequent sign-ins from this device. The trust period is refreshed each time the user signs in successfully.

### [Issuer](#issuer)

By adding an `issuer` you can set your application name for the 2fa application.

For example, if your user uses Google Auth, the default appName will show up as `Better Auth`. However, by using the following code, it will show up as `my-app-name`.

    twoFactor({
        issuer: "my-app-name"
    })

* * *

The plugin requires 1 additional fields in the `user` table and 1 additional table to store the two factor authentication data.

Table: `user`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| twoFactorEnabled | boolean |  | Whether two factor authentication is enabled for the user. |

Table: `twoFactor`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | The ID of the two factor authentication. |
| userId | string |  | The ID of the user |
| secret | string |  | The secret used to generate the TOTP code. |
| backupCodes | string |  | The backup codes used to recover access to the account if the user loses access to their phone or email. |

### [Server](#server)

**twoFactorTable**: The name of the table that stores the two factor authentication data. Default: `twoFactor`.

**skipVerificationOnEnable**: Skip the verification process before enabling two factor for a user.

**Issuer**: The issuer is the name of your application. It's used to generate TOTP codes. It'll be displayed in the authenticator apps.

**TOTP options**

these are options for TOTP.

**OTP options**

these are options for OTP.

**Backup Code Options**

backup codes are generated and stored in the database when the user enabled two factor authentication. This can be used to recover access to the account if the user loses access to their phone or email.

### [Client](#client)

To use the two factor plugin in the client, you need to add it on your plugins list.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { twoFactorClient } from "better-auth/client/plugins"
    
    const authClient =  createAuthClient({
        plugins: [
            twoFactorClient({ 
                onTwoFactorRedirect(){ 
                    window.location.href = "/2fa" // Handle the 2FA verification redirect
                } 
            }) 
        ]
    })

**Options**

`onTwoFactorRedirect`: A callback that will be called when the user needs to verify their 2FA code. This can be used to redirect the user to the 2FA page.</content>
</page>

<page>
  <title>API | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/api</url>
  <content>When you create a new Better Auth instance, it provides you with an `api` object. This object exposes every endpoint that exists in your Better Auth instance. And you can use this to interact with Better Auth server side.

Any endpoint added to Better Auth, whether from plugins or the core, will be accessible through the `api` object.

To call an API endpoint on the server, import your `auth` instance and call the endpoint using the `api` object.

server.ts

    import { betterAuth } from "better-auth";
    import { headers } from "next/headers";
    
    export const auth = betterAuth({
        //...
    })
    
    // calling get session on the server
    await auth.api.getSession({
        headers: await headers() // some endpoints might require headers
    })

### [Body, Headers, Query](#body-headers-query)

Unlike the client, the server needs the values to be passed as an object with the key `body` for the body, `headers` for the headers, and `query` for query parameters.

server.ts

    await auth.api.getSession({
        headers: await headers()
    })
    
    await auth.api.signInEmail({
        body: {
            email: "john@doe.com",
            password: "password"
        },
        headers: await headers() // optional but would be useful to get the user IP, user agent, etc.
    })
    
    await auth.api.verifyEmail({
        query: {
            token: "my_token"
        }
    })

Better Auth API endpoints are built on top of [better-call](https://github.com/bekacru/better-call), a tiny web framework that lets you call REST API endpoints as if they were regular functions and allows us to easily infer client types from the server.

### [Getting `headers` and `Response` Object](#getting-headers-and-response-object)

When you invoke an API endpoint on the server, it will return a standard JavaScript object or array directly as it's just a regular function call.

But there are times when you might want to get the `headers` or the `Response` object instead. For example, if you need to get the cookies or the headers.

#### [Getting `headers`](#getting-headers)

To get the `headers`, you can pass the `returnHeaders` option to the endpoint.

    const { headers, response } = await auth.api.signUpEmail({
    	returnHeaders: true,
    	body: {
    		email: "john@doe.com",
    		password: "password",
    		name: "John Doe",
    	},
    });

The `headers` will be a `Headers` object, which you can use to get the cookies or the headers.

    const cookies = headers.get("set-cookie");
    const headers = headers.get("x-custom-header");

#### [Getting `Response` Object](#getting-response-object)

To get the `Response` object, you can pass the `asResponse` option to the endpoint.

server.ts

    const response = await auth.api.signInEmail({
        body: {
            email: "",
            password: ""
        },
        asResponse: true
    })

### [Error Handling](#error-handling)

When you call an API endpoint on the server, it will throw an error if the request fails. You can catch the error and handle it as you see fit. The error instance is an instance of `APIError`.

server.ts

    import { APIError } from "better-auth/api";
    
    try {
        await auth.api.signInEmail({
            body: {
                email: "",
                password: ""
            }
        })
    } catch (error) {
        if (error instanceof APIError) {
            console.log(error.message, error.status)
        }
    }

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/concepts/api.mdx)</content>
</page>

<page>
  <title>Multi Session | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/multi-session</url>
  <content>The multi-session plugin allows users to maintain multiple active sessions across different accounts in the same browser. This plugin is useful for applications that require users to switch between multiple accounts without logging out.

### [Add the plugin to your **auth** config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { multiSession } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [ 
            multiSession(), 
        ] 
    })

### [Add the client Plugin](#add-the-client-plugin)

Add the client plugin and Specify where the user should be redirected if they need to verify 2nd factor

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { multiSessionClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            multiSessionClient()
        ]
    })

Whenever a user logs in, the plugin will add additional cookie to the browser. This cookie will be used to maintain multiple sessions across different accounts.

### [List all device sessions](#list-all-device-sessions)

To list all active sessions for the current user, you can call the `listDeviceSessions` method.

GET

/multi-session/list-device-sessions

    const { data, error } = await authClient.multiSession.listDeviceSessions();

### [Set active session](#set-active-session)

To set the active session, you can call the `setActive` method.

POST

/multi-session/set-active

    const { data, error } = await authClient.multiSession.setActive({    sessionToken: "some-session-token", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `sessionToken` | 
The session token to set as active.

 | `string` |

### [Revoke a session](#revoke-a-session)

To revoke a session, you can call the `revoke` method.

POST

/multi-session/revoke

    const { data, error } = await authClient.multiSession.revoke({    sessionToken: "some-session-token", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `sessionToken` | 
The session token to revoke.

 | `string` |

### [Signout and Revoke all sessions](#signout-and-revoke-all-sessions)

When a user logs out, the plugin will revoke all active sessions for the user. You can do this by calling the existing `signOut` method, which handles revoking all sessions automatically.

### [Max Sessions](#max-sessions)

You can specify the maximum number of sessions a user can have by passing the `maximumSessions` option to the plugin. By default, the plugin allows 5 sessions per device.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        plugins: [
            multiSession({
                maximumSessions: 3
            })
        ]
    })

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/multi-session.mdx)</content>
</page>

<page>
  <title>TypeScript | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/typescript</url>
  <content>Better Auth is designed to be type-safe. Both the client and server are built with TypeScript, allowing you to easily infer types.

### [Strict Mode](#strict-mode)

Better Auth is designed to work with TypeScript's strict mode. We recommend enabling strict mode in your TypeScript config file:

tsconfig.json

    {
      "compilerOptions": {
        "strict": true
      }
    }

if you can't set `strict` to `true`, you can enable `strictNullChecks`:

tsconfig.json

    {
      "compilerOptions": {
        "strictNullChecks": true,
      }
    }

If you're running into issues with TypeScript inference exceeding maximum length the compiler will serialize, then please make sure you're following the instructions above, as well as ensuring that both `declaration` and `composite` are not enabled.

Both the client SDK and the server offer types that can be inferred using the `$Infer` property. Plugins can extend base types like `User` and `Session`, and you can use `$Infer` to infer these types. Additionally, plugins can provide extra types that can also be inferred through `$Infer`.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    
    const authClient = createAuthClient()
    
    export type Session = typeof authClient.$Infer.Session

The `Session` type includes both `session` and `user` properties. The user property represents the user object type, and the `session` property represents the `session` object type.

You can also infer types on the server side.

auth.ts

    import { betterAuth } from "better-auth"
    import Database from "better-sqlite3"
    
    export const auth = betterAuth({
        database: new Database("database.db")
    })
    
    type Session = typeof auth.$Infer.Session

Better Auth allows you to add additional fields to the user and session objects. All additional fields are properly inferred and available on the server and client side.

    import { betterAuth } from "better-auth"
    import Database from "better-sqlite3"
    
    export const auth = betterAuth({
        database: new Database("database.db"),
        user: {
           additionalFields: {
              role: {
                  type: "string",
                  input: false
                } 
            }
        }
       
    })
    
    type Session = typeof auth.$Infer.Session

In the example above, we added a `role` field to the user object. This field is now available on the `Session` type.

### [The `input` property](#the-input-property)

The `input` property in an additional field configuration determines whether the field should be included in the user input. This property defaults to `true`, meaning the field will be part of the user input during operations like registration.

To prevent a field from being part of the user input, you must explicitly set `input: false`:

    additionalFields: {
        role: {
            type: "string",
            input: false
        }
    }

When `input` is set to `false`, the field will be excluded from user input, preventing users from passing a value for it.

By default, additional fields are included in the user input, which can lead to security vulnerabilities if not handled carefully. For fields that should not be set by the user, like a `role`, it is crucial to set `input: false` in the configuration.

### [Inferring Additional Fields on Client](#inferring-additional-fields-on-client)

To make sure proper type inference for additional fields on the client side, you need to inform the client about these fields. There are two approaches to achieve this, depending on your project structure:

1.  For Monorepo or Single-Project Setups

If your server and client code reside in the same project, you can use the `inferAdditionalFields` plugin to automatically infer the additional fields from your server configuration.

    import { inferAdditionalFields } from "better-auth/client/plugins";
    import { createAuthClient } from "better-auth/react";
    import type { auth } from "./auth";
    
    export const authClient = createAuthClient({
      plugins: [inferAdditionalFields<typeof auth>()],
    });

2.  For Separate Client-Server Projects

If your client and server are in separate projects, you'll need to manually specify the additional fields when creating the auth client.

    import { inferAdditionalFields } from "better-auth/client/plugins";
    
    export const authClient = createAuthClient({
      plugins: [inferAdditionalFields({
          user: {
            role: {
              type: "string"
            }
          }
      })],
    });

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/concepts/typescript.mdx)</content>
</page>

<page>
  <title>MySQL | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/mysql</url>
  <content>MySQL is a popular open-source relational database management system (RDBMS) that is widely used for building web applications and other types of software. It provides a flexible and scalable database solution that allows for efficient storage and retrieval of data. Read more here: [MySQL](https://www.mysql.com/).

Make sure you have MySQL installed and configured. Then, you can connect it straight into Better Auth.

auth.ts

    import { betterAuth } from "better-auth";
    import { createPool } from "mysql2/promise";
    
    export const auth = betterAuth({
      database: createPool({
        host: "localhost",
        user: "root",
        password: "password",
        database: "database",
        timezone: "Z", // Important to ensure consistent timezone values
      }),
    });

For more information, read Kysely's documentation to the [MySQLDialect](https://kysely-org.github.io/kysely-apidoc/classes/MysqlDialect.html).

The [Better Auth CLI](https://www.better-auth.com/docs/concepts/cli) allows you to generate or migrate your database schema based on your Better Auth configuration and plugins.

| 
MySQL Schema Generation

 | 

MySQL Schema Migration

 |
| --- | --- |
| ✅ Supported | ✅ Supported |

Database joins is useful when Better-Auth needs to fetch related data from multiple tables in a single query. Endpoints like `/get-session`, `/get-full-organization` and many others benefit greatly from this feature, seeing upwards of 2x to 3x performance improvements depending on database latency.

The Kysely MySQL dialect supports joins out of the box since version `1.4.0`.

To enable this feature, you need to set the `experimental.joins` option to `true` in your auth configuration.

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

It's possible that you may need to run migrations after enabling this feature.

MySQL is supported under the hood via the [Kysely](https://kysely.dev/) adapter, any database supported by Kysely would also be supported. ([Read more here](https://www.better-auth.com/docs/adapters/other-relational-databases))

If you're looking for performance improvements or tips, take a look at our guide to [performance optimizations](https://www.better-auth.com/docs/guides/optimizing-for-performance).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/adapters/mysql.mdx)</content>
</page>

<page>
  <title>Next.js Example | Better Auth</title>
  <url>https://www.better-auth.com/docs/examples/next-js</url>
  <content>This is an example of how to use Better Auth with Next.

**Implements the following features:** Email & Password . Social Sign-in . Passkeys . Email Verification . Password Reset . Two Factor Authentication . Profile Update . Session Management . Organization, Members and Roles

See [Demo](https://demo.better-auth.com/)

1.  Clone the code sandbox (or the repo) and open it in your code editor
2.  Move .env.example to .env and provide necessary variables
3.  Run the following commands
    
        pnpm install
        pnpm dev
    
4.  Open the browser and navigate to `http://localhost:3000`

### [SSO Login Example](#sso-login-example)

For this example, we utilize DummyIDP. Initiate the login from the [DummyIDP login](https://dummyidp.com/apps/app_01k16v4vb5yytywqjjvv2b3435/login), click "Proceed", and from here it will direct you to user's dashboard.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/examples/next-js.mdx)</content>
</page>

<page>
  <title>SQLite | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/sqlite</url>
  <content>SQLite is a lightweight, serverless, self-contained SQL database engine that is widely used for local data storage in applications. Read more [here.](https://www.sqlite.org/)

Better Auth supports multiple SQLite drivers. Choose the one that best fits your environment:

### [Better-SQLite3 (Recommended)](#better-sqlite3-recommended)

The most popular and stable SQLite driver for Node.js:

auth.ts

    import { betterAuth } from "better-auth";
    import Database from "better-sqlite3";
    
    export const auth = betterAuth({
      database: new Database("database.sqlite"),
    });

For more information, read Kysely's documentation to the [SqliteDialect](https://kysely-org.github.io/kysely-apidoc/classes/SqliteDialect.html).

### [Node.js Built-in SQLite (Experimental)](#nodejs-built-in-sqlite-experimental)

The `node:sqlite` module is still experimental and may change at any time. It requires Node.js 22.5.0 or later.

Starting from Node.js 22.5.0, you can use the built-in [SQLite](https://nodejs.org/api/sqlite.html) module:

auth.ts

    import { betterAuth } from "better-auth";
    import { DatabaseSync } from "node:sqlite";
    
    export const auth = betterAuth({
      database: new DatabaseSync("database.sqlite"),
    });

To run your application with Node.js SQLite:

    node your-app.js

### [Bun Built-in SQLite](#bun-built-in-sqlite)

You can also use the built-in [SQLite](https://bun.com/docs/api/sqlite) module in Bun, which is similar to the Node.js version:

auth.ts

    import { betterAuth } from "better-auth";
    import { Database } from "bun:sqlite";
    export const auth = betterAuth({
      database: new Database("database.sqlite"),
    });

The [Better Auth CLI](https://www.better-auth.com/docs/concepts/cli) allows you to generate or migrate your database schema based on your Better Auth configuration and plugins.

| 
SQLite Schema Generation

 | 

SQLite Schema Migration

 |
| --- | --- |
| ✅ Supported | ✅ Supported |

Database joins is useful when Better-Auth needs to fetch related data from multiple tables in a single query. Endpoints like `/get-session`, `/get-full-organization` and many others benefit greatly from this feature, seeing upwards of 2x to 3x performance improvements depending on database latency.

The Kysely SQLite dialect supports joins out of the box since version `1.4.0`. To enable this feature, you need to set the `experimental.joins` option to `true` in your auth configuration.

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

It's possible that you may need to run migrations after enabling this feature.

SQLite is supported under the hood via the [Kysely](https://kysely.dev/) adapter, any database supported by Kysely would also be supported. ([Read more here](https://www.better-auth.com/docs/adapters/other-relational-databases))

If you're looking for performance improvements or tips, take a look at our guide to [performance optimizations](https://www.better-auth.com/docs/guides/optimizing-for-performance).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/adapters/sqlite.mdx)</content>
</page>

<page>
  <title>PostgreSQL | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/postgresql</url>
  <content>PostgreSQL is a powerful, open-source relational database management system known for its advanced features, extensibility, and support for complex queries and large datasets. Read more [here](https://www.postgresql.org/).

Make sure you have PostgreSQL installed and configured. Then, you can connect it straight into Better Auth.

auth.ts

    import { betterAuth } from "better-auth";
    import { Pool } from "pg";
    
    export const auth = betterAuth({
      database: new Pool({
        connectionString: "postgres://user:password@localhost:5432/database",
      }),
    });

The [Better Auth CLI](https://www.better-auth.com/docs/concepts/cli) allows you to generate or migrate your database schema based on your Better Auth configuration and plugins.

| 
PostgreSQL Schema Generation

 | 

PostgreSQL Schema Migration

 |
| --- | --- |
| ✅ Supported | ✅ Supported |

Database joins is useful when Better-Auth needs to fetch related data from multiple tables in a single query. Endpoints like `/get-session`, `/get-full-organization` and many others benefit greatly from this feature, seeing upwards of 2x to 3x performance improvements depending on database latency.

The Kysely PostgreSQL dialect supports joins out of the box since version `1.4.0`. To enable this feature, you need to set the `experimental.joins` option to `true` in your auth configuration.

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

It's possible that you may need to run migrations after enabling this feature.

In most cases, the default schema is `public`. To have Better Auth use a non-default schema (e.g., `auth`) for its tables, you have several options:

### [Option 1: Set search\_path in connection string (Recommended)](#option-1-set-search_path-in-connection-string-recommended)

Append the `options` parameter to your connection URI:

auth.ts

    import { betterAuth } from "better-auth";
    import { Pool } from "pg";
    
    export const auth = betterAuth({
      database: new Pool({
        connectionString: "postgres://user:password@localhost:5432/database?options=-c search_path=auth",
      }),
    });

URL-encode if needed: `?options=-c%20search_path%3Dauth`.

### [Option 2: Set search\_path using Pool options](#option-2-set-search_path-using-pool-options)

auth.ts

    import { betterAuth } from "better-auth";
    import { Pool } from "pg";
    
    export const auth = betterAuth({
      database: new Pool({
        host: "localhost",
        port: 5432,
        user: "postgres",
        password: "password",
        database: "my-db",
        options: "-c search_path=auth",
      }),
    });

### [Option 3: Set default schema for database user](#option-3-set-default-schema-for-database-user)

Set the PostgreSQL user's default schema:

    ALTER USER your_user SET search_path TO auth;

After setting this, reconnect to apply the changes.

### [Prerequisites](#prerequisites)

Before using a non-default schema, ensure:

1.  **The schema exists:**
    
        CREATE SCHEMA IF NOT EXISTS auth;
    
2.  **The user has appropriate permissions:**
    
        GRANT ALL PRIVILEGES ON SCHEMA auth TO your_user;
        GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO your_user;
        ALTER DEFAULT PRIVILEGES IN SCHEMA auth GRANT ALL ON TABLES TO your_user;
    

### [How it works](#how-it-works)

The Better Auth CLI migration system automatically detects your configured `search_path`:

*   When running `npx @better-auth/cli migrate`, it inspects only the tables in your configured schema
*   Tables in other schemas (e.g., `public`) are ignored, preventing conflicts
*   All new tables are created in your specified schema

### [Troubleshooting](#troubleshooting)

**Issue:** "relation does not exist" error during migration

**Solution:** This usually means the schema doesn't exist or the user lacks permissions. Create the schema and grant permissions as shown above.

**Verifying your schema configuration:**

You can verify which schema Better Auth will use by checking the `search_path`:

    SHOW search_path;

This should return your custom schema (e.g., `auth`) as the first value.

PostgreSQL is supported under the hood via the [Kysely](https://kysely.dev/) adapter, any database supported by Kysely would also be supported. ([Read more here](https://www.better-auth.com/docs/adapters/other-relational-databases))

If you're looking for performance improvements or tips, take a look at our guide to [performance optimizations](https://www.better-auth.com/docs/guides/optimizing-for-performance).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/adapters/postgresql.mdx)</content>
</page>

<page>
  <title>Rate Limit | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/rate-limit</url>
  <content>Better Auth includes a built-in rate limiter to help manage traffic and prevent abuse. By default, in production mode, the rate limiter is set to:

*   Window: 60 seconds
*   Max Requests: 100 requests

Server-side requests made using `auth.api` aren't affected by rate limiting. Rate limits only apply to client-initiated requests.

You can easily customize these settings by passing the rateLimit object to the betterAuth function.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        rateLimit: {
            window: 10, // time window in seconds
            max: 100, // max requests in the window
        },
    })

Rate limiting is disabled in development mode by default. In order to enable it, set `enabled` to `true`:

auth.ts

    export const auth = betterAuth({
        rateLimit: {
            enabled: true,
            //...other options
        },
    })

In addition to the default settings, Better Auth provides custom rules for specific paths. For example:

*   `/sign-in/email`: Is limited to 3 requests within 10 seconds.

In addition, plugins also define custom rules for specific paths. For example, `twoFactor` plugin has custom rules:

*   `/two-factor/verify`: Is limited to 3 requests within 10 seconds.

These custom rules ensure that sensitive operations are protected with stricter limits.

### [Connecting IP Address](#connecting-ip-address)

Rate limiting uses the connecting IP address to track the number of requests made by a user. The default header checked is `x-forwarded-for`, which is commonly used in production environments. If you are using a different header to track the user's IP address, you'll need to specify it.

auth.ts

    export const auth = betterAuth({
        //...other options
        advanced: {
            ipAddress: {
              ipAddressHeaders: ["cf-connecting-ip"], // Cloudflare specific header example
          },
        },
        rateLimit: {
            enabled: true,
            window: 60, // time window in seconds
            max: 100, // max requests in the window
        },
    })

### [Rate Limit Window](#rate-limit-window)

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        //...other options
        rateLimit: {
            window: 60, // time window in seconds
            max: 100, // max requests in the window
        },
    })

You can also pass custom rules for specific paths.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        //...other options
        rateLimit: {
            window: 60, // time window in seconds
            max: 100, // max requests in the window
            customRules: {
                "/sign-in/email": {
                    window: 10,
                    max: 3,
                },
                "/two-factor/*": async (request)=> {
                    // custom function to return rate limit window and max
                    return {
                        window: 10,
                        max: 3,
                    }
                }
            },
        },
    })

If you like to disable rate limiting for a specific path, you can set it to `false` or return `false` from the custom rule function.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        //...other options
        rateLimit: {
            customRules: {
                "/get-session": false,
            },
        },
    })

### [Storage](#storage)

By default, rate limit data is stored in memory, which may not be suitable for many use cases, particularly in serverless environments. To address this, you can use a database, secondary storage, or custom storage for storing rate limit data.

**Using Database**

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        //...other options
        rateLimit: {
            storage: "database",
            modelName: "rateLimit", //optional by default "rateLimit" is used
        },
    })

Make sure to run `migrate` to create the rate limit table in your database.

    npx @better-auth/cli migrate

**Using Secondary Storage**

If a [Secondary Storage](https://www.better-auth.com/docs/concepts/database#secondary-storage) has been configured you can use that to store rate limit data.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        //...other options
        rateLimit: {
    		storage: "secondary-storage"
        },
    })

**Custom Storage**

If none of the above solutions suits your use case you can implement a `customStorage`.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        //...other options
        rateLimit: {
            customStorage: {
                get: async (key) => {
                    // get rate limit data
                },
                set: async (key, value) => {
                    // set rate limit data
                },
            },
        },
    })

When a request exceeds the rate limit, Better Auth returns the following header:

*   `X-Retry-After`: The number of seconds until the user can make another request.

To handle rate limit errors on the client side, you can manage them either globally or on a per-request basis. Since Better Auth clients wrap over Better Fetch, you can pass `fetchOptions` to handle rate limit errors

**Global Handling**

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    
    export const authClient = createAuthClient({
        fetchOptions: {
            onError: async (context) => {
                const { response } = context;
                if (response.status === 429) {
                    const retryAfter = response.headers.get("X-Retry-After");
                    console.log(`Rate limit exceeded. Retry after ${retryAfter} seconds`);
                }
            },
        }
    })

**Per Request Handling**

auth-client.ts

    import { authClient } from "./auth-client";
    
    await authClient.signIn.email({
        fetchOptions: {
            onError: async (context) => {
                const { response } = context;
                if (response.status === 429) {
                    const retryAfter = response.headers.get("X-Retry-After");
                    console.log(`Rate limit exceeded. Retry after ${retryAfter} seconds`);
                }
            },
        }
    })

### [Schema](#schema)

If you are using a database to store rate limit data you need this schema:

Table Name: `rateLimit`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Database ID |
| key | string | \- | Unique identifier for each rate limit key |
| count | integer | \- | Time window in seconds |
| lastRequest | bigint | \- | Max requests in the window |

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/concepts/rate-limit.mdx)</content>
</page>

<page>
  <title>Hooks | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/hooks#ctx</url>
  <content>Hooks in Better Auth let you "hook into" the lifecycle and execute custom logic. They provide a way to customize Better Auth's behavior without writing a full plugin.

We highly recommend using hooks if you need to make custom adjustments to an endpoint rather than making another endpoint outside of Better Auth.

**Before hooks** run _before_ an endpoint is executed. Use them to modify requests, pre validate data, or return early.

### [Example: Enforce Email Domain Restriction](#example-enforce-email-domain-restriction)

This hook ensures that users can only sign up if their email ends with `@example.com`:

auth.ts

    import { betterAuth } from "better-auth";
    import { createAuthMiddleware, APIError } from "better-auth/api";
    
    export const auth = betterAuth({
        hooks: {
            before: createAuthMiddleware(async (ctx) => {
                if (ctx.path !== "/sign-up/email") {
                    return;
                }
                if (!ctx.body?.email.endsWith("@example.com")) {
                    throw new APIError("BAD_REQUEST", {
                        message: "Email must end with @example.com",
                    });
                }
            }),
        },
    });

### [Example: Modify Request Context](#example-modify-request-context)

To adjust the request context before proceeding:

auth.ts

    import { betterAuth } from "better-auth";
    import { createAuthMiddleware } from "better-auth/api";
    
    export const auth = betterAuth({
        hooks: {
            before: createAuthMiddleware(async (ctx) => {
                if (ctx.path === "/sign-up/email") {
                    return {
                        context: {
                            ...ctx,
                            body: {
                                ...ctx.body,
                                name: "John Doe",
                            },
                        }
                    };
                }
            }),
        },
    });

**After hooks** run _after_ an endpoint is executed. Use them to modify responses.

### [Example: Send a notification to your channel when a new user is registered](#example-send-a-notification-to-your-channel-when-a-new-user-is-registered)

auth.ts

    import { betterAuth } from "better-auth";
    import { createAuthMiddleware } from "better-auth/api";
    import { sendMessage } from "@/lib/notification"
    
    export const auth = betterAuth({
        hooks: {
            after: createAuthMiddleware(async (ctx) => {
                if(ctx.path.startsWith("/sign-up")){
                    const newSession = ctx.context.newSession;
                    if(newSession){
                        sendMessage({
                            type: "user-register",
                            name: newSession.user.name,
                        })
                    }
                }
            }),
        },
    });

When you call `createAuthMiddleware` a `ctx` object is passed that provides a lot of useful properties. Including:

*   **Path:** `ctx.path` to get the current endpoint path.
*   **Body:** `ctx.body` for parsed request body (available for POST requests).
*   **Headers:** `ctx.headers` to access request headers.
*   **Request:** `ctx.request` to access the request object (may not exist in server-only endpoints).
*   **Query Parameters:** `ctx.query` to access query parameters.
*   **Context**: `ctx.context` auth related context, useful for accessing new session, auth cookies configuration, password hashing, config...

and more.

### [Request Response](#request-response)

This utilities allows you to get request information and to send response from a hook.

#### [JSON Responses](#json-responses)

Use `ctx.json` to send JSON responses:

    const hook = createAuthMiddleware(async (ctx) => {
        return ctx.json({
            message: "Hello World",
        });
    });

#### [Redirects](#redirects)

Use `ctx.redirect` to redirect users:

    import { createAuthMiddleware } from "better-auth/api";
    
    const hook = createAuthMiddleware(async (ctx) => {
        throw ctx.redirect("/sign-up/name");
    });

#### [Cookies](#cookies)

*   Set cookies: `ctx.setCookies` or `ctx.setSignedCookie`.
*   Get cookies: `ctx.getCookies` or `ctx.getSignedCookie`.

Example:

    import { createAuthMiddleware } from "better-auth/api";
    
    const hook = createAuthMiddleware(async (ctx) => {
        ctx.setCookies("my-cookie", "value");
        await ctx.setSignedCookie("my-signed-cookie", "value", ctx.context.secret, {
            maxAge: 1000,
        });
    
        const cookie = ctx.getCookies("my-cookie");
        const signedCookie = await ctx.getSignedCookie("my-signed-cookie");
    });

#### [Errors](#errors)

Throw errors with `APIError` for a specific status code and message:

    import { createAuthMiddleware, APIError } from "better-auth/api";
    
    const hook = createAuthMiddleware(async (ctx) => {
        throw new APIError("BAD_REQUEST", {
            message: "Invalid request",
        });
    });

### [Context](#context)

The `ctx` object contains another `context` object inside that's meant to hold contexts related to auth. Including a newly created session on after hook, cookies configuration, password hasher and so on.

#### [New Session](#new-session)

The newly created session after an endpoint is run. This only exist in after hook.

auth.ts

    createAuthMiddleware(async (ctx) => {
        const newSession = ctx.context.newSession
    });

#### [Returned](#returned)

The returned value from the hook is passed to the next hook in the chain.

auth.ts

    createAuthMiddleware(async (ctx) => {
        const returned = ctx.context.returned; //this could be a successful response or an APIError
    });

#### [Response Headers](#response-headers)

The response headers added by endpoints and hooks that run before this hook.

auth.ts

    createAuthMiddleware(async (ctx) => {
        const responseHeaders = ctx.context.responseHeaders;
    });

#### [Predefined Auth Cookies](#predefined-auth-cookies)

Access BetterAuth’s predefined cookie properties:

auth.ts

    createAuthMiddleware(async (ctx) => {
        const cookieName = ctx.context.authCookies.sessionToken.name;
    });

#### [Secret](#secret)

You can access the `secret` for your auth instance on `ctx.context.secret`

#### [Password](#password)

The password object provider `hash` and `verify`

*   `ctx.context.password.hash`: let's you hash a given password.
*   `ctx.context.password.verify`: let's you verify given `password` and a `hash`.

#### [Adapter](#adapter)

Adapter exposes the adapter methods used by Better Auth. Including `findOne`, `findMany`, `create`, `delete`, `update` and `updateMany`. You generally should use your actually `db` instance from your orm rather than this adapter.

#### [Internal Adapter](#internal-adapter)

These are calls to your db that perform specific actions. `createUser`, `createSession`, `updateSession`...

This may be useful to use instead of using your db directly to get access to `databaseHooks`, proper `secondaryStorage` support and so on. If you're make a query similar to what exist in this internal adapter actions it's worth a look.

#### [generateId](#generateid)

You can use `ctx.context.generateId` to generate Id for various reasons.

If you need to reuse a hook across multiple endpoints, consider creating a plugin. Learn more in the [Plugins Documentation](https://www.better-auth.com/docs/concepts/plugins).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/concepts/hooks.mdx)</content>
</page>

<page>
  <title>Client | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/client</url>
  <content>Better Auth offers a client library compatible with popular frontend frameworks like React, Vue, Svelte, and more. This client library includes a set of functions for interacting with the Better Auth server. Each framework's client library is built on top of a core client library that is framework-agnostic, so that all methods and hooks are consistently available across all client libraries.

If you haven't already, install better-auth.

Import `createAuthClient` from the package for your framework (e.g., "better-auth/react" for React). Call the function to create your client. Pass the base URL of your auth server. If the auth server is running on the same domain as your client, you can skip this step.

If you're using a different base path other than `/api/auth`, make sure to pass the whole URL, including the path. (e.g., `http://localhost:3000/custom-path/auth`)

Once you've created your client instance, you can use the client to interact with the Better Auth server. The client provides a set of functions by default and they can be extended with plugins.

**Example: Sign In**

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    const authClient = createAuthClient()
    
    await authClient.signIn.email({
        email: "test@user.com",
        password: "password1234"
    })

### [Hooks](#hooks)

In addition to the standard methods, the client provides hooks to easily access different reactive data. Every hook is available in the root object of the client and they all start with `use`.

**Example: useSession**

### [Fetch Options](#fetch-options)

The client uses a library called [better fetch](https://better-fetch.vercel.app/) to make requests to the server.

Better fetch is a wrapper around the native fetch API that provides a more convenient way to make requests. It's created by the same team behind Better Auth and is designed to work seamlessly with it.

You can pass any default fetch options to the client by passing `fetchOptions` object to the `createAuthClient`.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    
    const authClient = createAuthClient({
        fetchOptions: {
            //any better-fetch options
        },
    })

You can also pass fetch options to most of the client functions. Either as the second argument or as a property in the object.

auth-client.ts

    await authClient.signIn.email({
        email: "email@email.com",
        password: "password1234",
    }, {
        onSuccess(ctx) {
                //      
        }
    })
    
    //or
    
    await authClient.signIn.email({
        email: "email@email.com",
        password: "password1234",
        fetchOptions: {
            onSuccess(ctx) {
                //      
            }
        },
    })

### [Disabling Hook Rerenders](#disabling-hook-rerenders)

Certain endpoints, upon successful response, will trigger atom signals and cause hooks like `useSession` to rerender. This is useful for keeping your UI in sync with authentication state changes.

However, there are cases where you might want to make an endpoint call without triggering hook rerenders. For example, when updating user preferences that don't affect the session, or when you want to manually control when hooks update.

You can disable hook rerenders for a specific endpoint call by setting `disableSignal: true` in the fetch options:

auth-client.ts

    // As the second argument
    await authClient.updateUser({
        name: "New Name"
    }, {
        disableSignal: true
    })
    
    // Or within fetchOptions
    await authClient.updateUser({
        name: "New Name",
        fetchOptions: {
            disableSignal: true
        }
    })

When `disableSignal` is set to `true`, the endpoint call will complete successfully, but hooks like `useSession` won't automatically rerender. You can manually trigger a refetch if needed:

auth-client.ts

    const { refetch } = authClient.useSession()
    
    await authClient.updateUser({
        name: "New Name"
    }, {
        disableSignal: true,
        onSuccess() {
            // Manually refetch session if needed
            refetch()
        }
    })

### [Handling Errors](#handling-errors)

Most of the client functions return a response object with the following properties:

*   `data`: The response data.
*   `error`: The error object if there was an error.

The error object contains the following properties:

*   `message`: The error message. (e.g., "Invalid email or password")
*   `status`: The HTTP status code.
*   `statusText`: The HTTP status text.

auth-client.ts

    const { data, error } = await authClient.signIn.email({
        email: "email@email.com",
        password: "password1234"
    })
    if (error) {
        //handle error
    }

If the action accepts a `fetchOptions` option, you can pass an `onError` callback to handle errors.

auth-client.ts

    
    await authClient.signIn.email({
        email: "email@email.com",
        password: "password1234",
    }, {
        onError(ctx) {
            //handle error
        }
    })
    
    //or
    await authClient.signIn.email({
        email: "email@email.com",
        password: "password1234",
        fetchOptions: {
            onError(ctx) {
                //handle error
            }
        }
    })

Hooks like `useSession` also return an error object if there was an error fetching the session. On top of that, they also return an `isPending` property to indicate if the request is still pending.

auth-client.ts

    const { data, error, isPending } = useSession()
    if (error) {
        //handle error
    }

#### [Error Codes](#error-codes)

The client instance contains $ERROR\_CODES object that contains all the error codes returned by the server. You can use this to handle error translations or custom error messages.

auth-client.ts

    const authClient = createAuthClient();
    
    type ErrorTypes = Partial<
    	Record<
    		keyof typeof authClient.$ERROR_CODES,
    		{
    			en: string;
    			es: string;
    		}
    	>
    >;
    
    const errorCodes = {
    	USER_ALREADY_EXISTS: {
    		en: "user already registered",
    		es: "usuario ya registrado",
    	},
    } satisfies ErrorTypes;
    
    const getErrorMessage = (code: string, lang: "en" | "es") => {
    	if (code in errorCodes) {
    		return errorCodes[code as keyof typeof errorCodes][lang];
    	}
    	return "";
    };
    
    
    const { error } = await authClient.signUp.email({
    	email: "user@email.com",
    	password: "password",
    	name: "User",
    });
    if(error?.code){
        alert(getErrorMessage(error.code, "en"));
    }

### [Plugins](#plugins)

You can extend the client with plugins to add more functionality. Plugins can add new functions to the client or modify existing ones.

**Example: Magic Link Plugin**

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { magicLinkClient } from "better-auth/client/plugins"
    
    const authClient = createAuthClient({
        plugins: [
            magicLinkClient()
        ]
    })

once you've added the plugin, you can use the new functions provided by the plugin.

auth-client.ts

    await authClient.signIn.magicLink({
        email: "test@email.com"
    })</content>
</page>

<page>
  <title>Drizzle ORM Adapter | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/drizzle#joins-experimental</url>
  <content>Drizzle ORM is a powerful and flexible ORM for Node.js and TypeScript. It provides a simple and intuitive API for working with databases, and supports a wide range of databases including MySQL, PostgreSQL, SQLite, and more.

Before getting started, make sure you have Drizzle installed and configured. For more information, see [Drizzle Documentation](https://orm.drizzle.team/docs/overview/)

You can use the Drizzle adapter to connect to your database as follows.

auth.ts

    import { betterAuth } from "better-auth";
    import { drizzleAdapter } from "better-auth/adapters/drizzle";
    import { db } from "./database.ts";
    
    export const auth = betterAuth({
      database: drizzleAdapter(db, {
        provider: "sqlite", // or "pg" or "mysql"
      }), 
      //... the rest of your config
    });

The [Better Auth CLI](https://www.better-auth.com/docs/concepts/cli) allows you to generate or migrate your database schema based on your Better Auth configuration and plugins.

To generate the schema required by Better Auth, run the following command:

To generate and apply the migration, run the following commands:

Database joins is useful when Better-Auth needs to fetch related data from multiple tables in a single query. Endpoints like `/get-session`, `/get-full-organization` and many others benefit greatly from this feature, seeing upwards of 2x to 3x performance improvements depending on database latency.

The Drizzle adapter supports joins out of the box since version `1.4.0`. To enable this feature, you need to set the `experimental.joins` option to `true` in your auth configuration.

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

Please make sure that your Drizzle schema has the necessary relations defined. If you do not see any relations in your Drizzle schema, you can manually add them using the [`relation`](https://orm.drizzle.team/docs/relations) drizzle-orm function or run our latest CLI version `npx @better-auth/cli@latest generate` to generate a new Drizzle schema with the relations.

Additionally, you're required to pass each [relation](https://orm.drizzle.team/docs/relations) through the drizzle adapter schema object.

The Drizzle adapter expects the schema you define to match the table names. For example, if your Drizzle schema maps the `user` table to `users`, you need to manually pass the schema and map it to the user table.

    import { betterAuth } from "better-auth";
    import { db } from "./drizzle";
    import { drizzleAdapter } from "better-auth/adapters/drizzle";
    import { schema } from "./schema";
    
    export const auth = betterAuth({
      database: drizzleAdapter(db, {
        provider: "sqlite", // or "pg" or "mysql"
        schema: {
          ...schema,
          user: schema.users,
        },
      }),
    });

You can either modify the provided schema values like the example above, or you can mutate the auth config's `modelName` property directly. For example:

    export const auth = betterAuth({
      database: drizzleAdapter(db, {
        provider: "sqlite", // or "pg" or "mysql"
        schema,
      }),
      user: {
        modelName: "users", 
      }
    });

We map field names based on property you passed to your Drizzle schema. For example, if you want to modify the `email` field to `email_address`, you simply need to change the Drizzle schema to:

    export const user = mysqlTable("user", {
      // Changed field name without changing the schema property name
      // This allows drizzle & better-auth to still use the original field name,
      // while your DB uses the modified field name
      email: varchar("email_address", { length: 255 }).notNull().unique(), 
      // ... others
    });

You can either modify the Drizzle schema like the example above, or you can mutate the auth config's `fields` property directly. For example:

    export const auth = betterAuth({
      database: drizzleAdapter(db, {
        provider: "sqlite", // or "pg" or "mysql"
        schema,
      }),
      user: {
        fields: {
          email: "email_address", 
        }
      }
    });

If all your tables are using plural form, you can just pass the `usePlural` option:

    export const auth = betterAuth({
      database: drizzleAdapter(db, {
        ...
        usePlural: true,
      }),
    });

*   If you're looking for performance improvements or tips, take a look at our guide to [performance optimizations](https://www.better-auth.com/docs/guides/optimizing-for-performance).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/adapters/drizzle.mdx)</content>
</page>

<page>
  <title>Remix Example | Better Auth</title>
  <url>https://www.better-auth.com/docs/examples/remix</url>
  <content>This is an example of how to use Better Auth with Remix.

**Implements the following features:** Email & Password . Social Sign-in with Google . Passkeys . Email Verification . Password Reset . Two Factor Authentication . Profile Update . Session Management

[How to run](#how-to-run)
-------------------------

1.  Clone the code sandbox (or the repo) and open it in your code editor
2.  Provide .env file with by copying the `.env.example` file and adding the variables
3.  Run the following commands
    
        pnpm install
        pnpm run dev
    
4.  Open the browser and navigate to `http://localhost:3000`</content>
</page>

<page>
  <title>Prisma | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/prisma#joins-experimental</url>
  <content>Prisma ORM is an open-source database toolkit that simplifies database access and management in applications by providing a type-safe query builder and an intuitive data modeling interface.

Before getting started, make sure you have Prisma installed and configured. For more information, see [Prisma Documentation](https://www.prisma.io/docs/)

You can use the Prisma adapter to connect to your database as follows.

auth.ts

    import { betterAuth } from "better-auth";
    import { prismaAdapter } from "better-auth/adapters/prisma";
    import { PrismaClient } from "@prisma/client";
    
    const prisma = new PrismaClient();
    
    export const auth = betterAuth({
      database: prismaAdapter(prisma, {
        provider: "sqlite",
      }),
    });

Starting from Prisma 7, the `output` path field is required. If you have configured a custom output path in your `schema.prisma` file (e.g., `output = "../src/generated/prisma"`), make sure to import the Prisma client from that location instead of `@prisma/client`. For more information, see [here](https://www.prisma.io/docs/orm/prisma-client/setup-and-configuration/generating-prisma-client#the-location-of-prisma-client).

The [Better Auth CLI](https://www.better-auth.com/docs/concepts/cli) allows you to generate or migrate your database schema based on your Better Auth configuration and plugins.

| 
Prisma Schema Generation

 | 

Prisma Schema Migration

 |
| --- | --- |
| ✅ Supported | ❌ Not Supported |

Database joins is useful when Better-Auth needs to fetch related data from multiple tables in a single query. Endpoints like `/get-session`, `/get-full-organization` and many others benefit greatly from this feature, seeing upwards of 2x to 3x performance improvements depending on database latency.

The Prisma adapter supports joins out of the box since version `1.4.0`. To enable this feature, you need to set the `experimental.joins` option to `true` in your auth configuration.

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

Please make sure that your Prisma schema has the necessary relations defined. If you do not see any relations in your Prisma schema, you can manually add them using the `@relation` directive or run our latest CLI version `npx @better-auth/cli@latest generate` to generate a new Prisma schema with the relations.

*   If you're looking for performance improvements or tips, take a look at our guide to [performance optimizations](https://www.better-auth.com/docs/guides/optimizing-for-performance).
*   [How to use Prisma ORM with Better Auth and Next.js](https://www.prisma.io/docs/guides/betterauth-nextjs)

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/adapters/prisma.mdx)</content>
</page>

<page>
  <title>Email | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/email</url>
  <content>Email is a key part of Better Auth, required for all users regardless of their authentication method. Better Auth provides email and password authentication out of the box, and a lot of utilities to help you manage email verification, password reset, and more.

Email verification is a security feature that ensures users provide a valid email address. It helps prevent spam and abuse by confirming that the email address belongs to the user. In this guide, you'll get a walk through of how to implement token based email verification in your app. To use otp based email verification, check out the [OTP Verification](https://www.better-auth.com/docs/plugins/email-otp) guide.

### [Adding Email Verification to Your App](#adding-email-verification-to-your-app)

To enable email verification, you need to pass a function that sends a verification email with a link.

*   **sendVerificationEmail**: This function is triggered when email verification starts. It accepts a data object with the following properties:
    *   `user`: The user object containing the email address.
    *   `url`: The verification URL the user must click to verify their email.
    *   `token`: The verification token used to complete the email verification to be used when implementing a custom verification URL.

and a `request` object as the second parameter.

auth.ts

    import { betterAuth } from 'better-auth';
    import { sendEmail } from './email'; // your email sending function
    
    export const auth = betterAuth({
        emailVerification: {
            sendVerificationEmail: async ({ user, url, token }, request) => {
                void sendEmail({
                    to: user.email,
                    subject: 'Verify your email address',
                    text: `Click the link to verify your email: ${url}`
                })
            }
        }
    })

Avoid awaiting the email sending to prevent timing attacks. On serverless platforms, use `waitUntil` or similar to ensure the email is sent.

### [Triggering Email Verification](#triggering-email-verification)

You can initiate email verification in several ways:

#### [1\. During Sign-up](#1-during-sign-up)

To automatically send a verification email at signup, set `emailVerification.sendOnSignUp` to `true`.

auth.ts

    import { betterAuth } from 'better-auth';
    
    export const auth = betterAuth({
        emailVerification: {
            sendOnSignUp: true
        }
    })

This sends a verification email when a user signs up. For social logins, email verification status is read from the SSO.

With `sendOnSignUp` enabled, when the user logs in with an SSO that does not claim the email as verified, Better Auth will dispatch a verification email, but the verification is not required to login even when `requireEmailVerification` is enabled.

#### [2\. Require Email Verification](#2-require-email-verification)

If you enable require email verification, users must verify their email before they can log in. And every time a user tries to sign in, `sendVerificationEmail` is called.

This only works if you have `sendVerificationEmail` implemented, if `sendOnSignIn` is set to true and if the user is trying to sign in with email and password.

auth.ts

    export const auth = betterAuth({
      emailVerification: {
        sendVerificationEmail: async ({ user, url }) => {
          void sendEmail({
            to: user.email,
            subject: "Verify your email address",
            text: `Click the link to verify your email: ${url}`,
          });
        },
        sendOnSignIn: true,
      },
      emailAndPassword: {
        requireEmailVerification: true,
      },
    });

If a user tries to sign in without verifying their email, you can handle the error and show a message to the user.

auth-client.ts

    await authClient.signIn.email({
        email: "email@example.com",
        password: "password"
    }, {
        onError: (ctx) => {
            // Handle the error
            if(ctx.error.status === 403) {
                alert("Please verify your email address")
            }
            //you can also show the original error message
            alert(ctx.error.message)
        }
    })

#### [3\. Manually](#3-manually)

You can also manually trigger email verification by calling `sendVerificationEmail`.

    await authClient.sendVerificationEmail({
        email: "user@email.com",
        callbackURL: "/" // The redirect URL after verification
    })

### [Verifying the Email](#verifying-the-email)

If the user clicks the provided verification URL, their email is automatically verified, and they are redirected to the `callbackURL`.

For manual verification, you can send the user a custom link with the `token` and call the `verifyEmail` function.

    await authClient.verifyEmail({
        query: {
            token: "" // Pass the token here
        }
    })

### [Auto Sign In After Verification](#auto-sign-in-after-verification)

To sign in the user automatically after they successfully verify their email, set the `autoSignInAfterVerification` option to `true`:

    const auth = betterAuth({
        //...your other options
        emailVerification: {
            autoSignInAfterVerification: true
        }
    })

### [Callback after successful email verification](#callback-after-successful-email-verification)

You can run custom code immediately after a user verifies their email using the `afterEmailVerification` callback. This is useful for any side-effects you want to trigger, like granting access to special features or logging the event.

The `afterEmailVerification` function runs automatically when a user's email is confirmed, receiving the `user` object and `request` details so you can perform actions for that specific user.

Here's how you can set it up:

auth.ts

    import { betterAuth } from 'better-auth';
    
    export const auth = betterAuth({
        emailVerification: {
            async afterEmailVerification(user, request) {
                // Your custom logic here, e.g., grant access to premium features
                console.log(`${user.email} has been successfully verified!`);
            }
        }
    })

Password reset allows users to reset their password if they forget it. Better Auth provides a simple way to implement password reset functionality.

You can enable password reset by passing a function that sends a password reset email with a link.

auth.ts

    import { betterAuth } from 'better-auth';
    import { sendEmail } from './email'; // your email sending function
    
    export const auth = betterAuth({
        emailAndPassword: {
            enabled: true,
            sendResetPassword: async ({ user, url, token }, request) => {
                void sendEmail({
                    to: user.email,
                    subject: 'Reset your password',
                    text: `Click the link to reset your password: ${url}`
                })
            }
        }
    })

Avoid awaiting the email sending to prevent timing attacks. On serverless platforms, use `waitUntil` or similar to ensure the email is sent.

Check out the [Email and Password](https://www.better-auth.com/docs/authentication/email-password#forget-password) guide for more details on how to implement password reset in your app. Also you can check out the [Otp verification](https://www.better-auth.com/docs/plugins/email-otp#reset-password) guide for how to implement password reset with OTP in your app.</content>
</page>

<page>
  <title>Cookies | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/cookies</url>
  <content>Cookies are used to store data such as session tokens, session data, OAuth state, and more. All cookies are signed using the `secret` key provided in the auth options or the `BETTER_AUTH_SECRET` environment variable.

### [Cookie Prefix](#cookie-prefix)

By default, Better Auth cookies follow the format `${prefix}.${cookie_name}`. The default prefix is "better-auth". You can change the prefix by setting `cookiePrefix` in the `advanced` object of the auth options.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        advanced: {
            cookiePrefix: "my-app"
        }
    })

### [Custom Cookies](#custom-cookies)

All cookies are `httpOnly` and `secure` when the server is running in production mode.

If you want to set custom cookie names and attributes, you can do so by setting `cookieOptions` in the `advanced` object of the auth options.

By default, Better Auth uses the following cookies:

*   `session_token` to store the session token
*   `session_data` to store the session data if cookie cache is enabled
*   `dont_remember` to store the flag when `rememberMe` is disabled

Plugins may also use cookies to store data. For example, the Two Factor Authentication plugin uses the `two_factor` cookie to store the two-factor authentication state.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        advanced: {
            cookies: {
                session_token: {
                    name: "custom_session_token",
                    attributes: {
                        // Set custom cookie attributes
                    }
                },
            }
        }
    })

### [Cross Subdomain Cookies](#cross-subdomain-cookies)

Sometimes you may need to share cookies across subdomains. For example, if you authenticate on `auth.example.com`, you may also want to access the same session on `app.example.com`.

The `domain` attribute controls which domains can access the cookie. Setting it to your root domain (e.g. `example.com`) makes the cookie accessible across all subdomains. For security, follow these guidelines:

1.  Only enable cross-subdomain cookies if it's necessary
2.  Set the domain to the most specific scope needed (e.g. `app.example.com` instead of `.example.com`)
3.  Be cautious of untrusted subdomains that could potentially access these cookies
4.  Consider using separate domains for untrusted services (e.g. `status.company.com` vs `app.company.com`)

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        advanced: {
            crossSubDomainCookies: {
                enabled: true,
                domain: "app.example.com", // your domain
            },
        },
        trustedOrigins: [
            'https://example.com',
            'https://app1.example.com',
            'https://app2.example.com',
        ],
    })

### [Secure Cookies](#secure-cookies)

By default, cookies are secure only when the server is running in production mode. You can force cookies to be always secure by setting `useSecureCookies` to `true` in the `advanced` object in the auth options.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        advanced: {
            useSecureCookies: true
        }
    })

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/concepts/cookies.mdx)</content>
</page>

<page>
  <title>MS SQL | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/mssql</url>
  <content>Microsoft SQL Server is a relational database management system developed by Microsoft, designed for enterprise-level data storage, management, and analytics with robust security and scalability features. Read more [here](https://en.wikipedia.org/wiki/Microsoft_SQL_Server).

[Example Usage](#example-usage)
-------------------------------

Make sure you have MS SQL installed and configured. Then, you can connect it straight into Better Auth.

auth.ts

    import { betterAuth } from "better-auth";
    import { MssqlDialect } from "kysely";
    import * as Tedious from 'tedious'
    import * as Tarn from 'tarn'
    
    const dialect = new MssqlDialect({
      tarn: {
        ...Tarn,
        options: {
          min: 0,
          max: 10,
        },
      },
      tedious: {
        ...Tedious,
        connectionFactory: () => new Tedious.Connection({
          authentication: {
            options: {
              password: 'password',
              userName: 'username',
            },
            type: 'default',
          },
          options: {
            database: 'some_db',
            port: 1433,
            trustServerCertificate: true,
          },
          server: 'localhost',
        }),
      },
      TYPES: {
    		...Tedious.TYPES,
    		DateTime: Tedious.TYPES.DateTime2,
    	},
    })
    
    export const auth = betterAuth({
      database: {
        dialect,
        type: "mssql"
      }
    });
    

For more information, read Kysely's documentation to the [MssqlDialect](https://kysely-org.github.io/kysely-apidoc/classes/MssqlDialect.html).

[Schema generation & migration](#schema-generation--migration)
--------------------------------------------------------------

The [Better Auth CLI](https://www.better-auth.com/docs/concepts/cli) allows you to generate or migrate your database schema based on your Better Auth configuration and plugins.

| 
MS SQL Schema Generation

 | 

MS SQL Schema Migration

 |
| --- | --- |
| ✅ Supported | ✅ Supported |

Schema Generation

    npx @better-auth/cli@latest generate

Schema Migration

    npx @better-auth/cli@latest migrate

[Joins (Experimental)](#joins-experimental)
-------------------------------------------

Database joins is useful when Better-Auth needs to fetch related data from multiple tables in a single query. Endpoints like `/get-session`, `/get-full-organization` and many others benefit greatly from this feature, seeing upwards of 2x to 3x performance improvements depending on database latency.

The Kysely MS SQL dialect supports joins out of the box since version `1.4.0`. To enable this feature, you need to set the `experimental.joins` option to `true` in your auth configuration.

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

It's possible that you may need to run migrations after enabling this feature.

[Additional Information](#additional-information)
-------------------------------------------------

MS SQL is supported under the hood via the [Kysely](https://kysely.dev/) adapter, any database supported by Kysely would also be supported. ([Read more here](https://www.better-auth.com/docs/adapters/other-relational-databases))

If you're looking for performance improvements or tips, take a look at our guide to [performance optimizations](https://www.better-auth.com/docs/guides/optimizing-for-performance).</content>
</page>

<page>
  <title>Expo Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/expo</url>
  <content>Expo is a popular framework for building cross-platform apps with React Native. Better Auth supports both Expo native and web apps.

Before using Better Auth with Expo, make sure you have a Better Auth backend set up. You can either use a separate server or leverage Expo's new [API Routes](https://docs.expo.dev/router/reference/api-routes) feature to host your Better Auth instance.

To get started, check out our [installation](https://www.better-auth.com/docs/installation) guide for setting up Better Auth on your server. If you prefer to check out the full example, you can find it [here](https://github.com/better-auth/examples/tree/main/expo-example).

To use the new API routes feature in Expo to host your Better Auth instance you can create a new API route in your Expo app and mount the Better Auth handler.

app/api/auth/\[...auth\]+api.ts

    import { auth } from "@/lib/auth"; // import Better Auth handler
    
    const handler = auth.handler;
    export { handler as GET, handler as POST }; // export handler for both GET and POST requests

Install both the Better Auth package and Expo plugin into your server application.

You also need to install both the Better Auth package and Expo plugin into your Expo application.

If you plan on using our social integrations (Google, Apple etc.) then there are a few more dependencies that are required in your Expo app. In the default Expo template these are already installed so you may be able to skip this step if you have these dependencies already.

Add the Expo plugin to your Better Auth server.

lib/auth.ts

    import { betterAuth } from "better-auth";
    import { expo } from "@better-auth/expo";
    
    export const auth = betterAuth({
        plugins: [expo()],
        emailAndPassword: { 
            enabled: true, // Enable authentication using email and password.
          }, 
    });

To initialize Better Auth in your Expo app, you need to call `createAuthClient` with the base URL of your Better Auth backend. Make sure to import the client from `/react`.

Make sure you install the `expo-secure-store` package into your Expo app. This is used to store the session data and cookies securely.

You need to also import client plugin from `@better-auth/expo/client` and pass it to the `plugins` array when initializing the auth client.

This is important because:

*   **Social Authentication Support:** enables social auth flows by handling authorization URLs and callbacks within the Expo web browser.
*   **Secure Cookie Management:** stores cookies securely and automatically adds them to the headers of your auth requests.

lib/auth-client.ts

    import { createAuthClient } from "better-auth/react";
    import { expoClient } from "@better-auth/expo/client";
    import * as SecureStore from "expo-secure-store";
    
    export const authClient = createAuthClient({
        baseURL: "http://localhost:8081", // Base URL of your Better Auth backend.
        plugins: [
            expoClient({
                scheme: "myapp",
                storagePrefix: "myapp",
                storage: SecureStore,
            })
        ]
    });

Be sure to include the full URL, including the path, if you've changed the default path from `/api/auth`.

Better Auth uses deep links to redirect users back to your app after authentication. To enable this, you need to add your app's scheme to the `trustedOrigins` list in your Better Auth config.

First, make sure you have a scheme defined in your `app.json` file.

app.json

    {
        "expo": {
            "scheme": "myapp"
        }
    }

Then, update your Better Auth config to include the scheme in the `trustedOrigins` list.

auth.ts

    export const auth = betterAuth({
        trustedOrigins: ["myapp://"]
    })

If you have multiple schemes or need to support deep linking with various paths, you can use specific patterns or wildcards:

auth.ts

    export const auth = betterAuth({
        trustedOrigins: [
            // Basic scheme
            "myapp://", 
            
            // Production & staging schemes
            "myapp-prod://",
            "myapp-staging://",
            
            // Wildcard support for all paths following the scheme
            "myapp://*"
        ]
    })

### [Development Mode](#development-mode)

During development, Expo uses the `exp://` scheme with your device's local IP address. To support this, you can use wildcards to match common local IP ranges:

auth.ts

    export const auth = betterAuth({
        trustedOrigins: [
            "myapp://",
            
            // Development mode - Expo's exp:// scheme with local IP ranges
            ...(process.env.NODE_ENV === "development" ? [
                "exp://*/*",                 // Trust all Expo development URLs
                "exp://10.0.0.*:*/*",        // Trust 10.0.0.x IP range
                "exp://192.168.*.*:*/*",     // Trust 192.168.x.x IP range
                "exp://172.*.*.*:*/*",       // Trust 172.x.x.x IP range
                "exp://localhost:*/*"        // Trust localhost
            ] : [])
        ]
    })

The wildcard patterns for `exp://` should only be used in development. In production, use your app's specific scheme (e.g., `myapp://`).

To resolve Better Auth exports you'll need to enable `unstable_enablePackageExports` in your metro config.

metro.config.js

    const { getDefaultConfig } = require("expo/metro-config");
    
    const config = getDefaultConfig(__dirname)
    
    config.resolver.unstable_enablePackageExports = true; 
    
    module.exports = config;

In case you don't have a `metro.config.js` file in your project run `npx expo customize metro.config.js`.

If you can't enable `unstable_enablePackageExports` option, you can use [babel-plugin-module-resolver](https://github.com/tleunen/babel-plugin-module-resolver) to manually resolve the paths.

babel.config.js

    module.exports = function (api) {
        api.cache(true);
        return {
            presets: ["babel-preset-expo"],
            plugins: [
                [
                    "module-resolver",
                    {
                        alias: {
                            "better-auth/react": "./node_modules/better-auth/dist/client/react/index.cjs",
                            "better-auth/client/plugins": "./node_modules/better-auth/dist/client/plugins/index.cjs",
                            "@better-auth/expo/client": "./node_modules/@better-auth/expo/dist/client.cjs",
                        },
                    },
                ],
            ],
        }
    }

In case you don't have a `babel.config.js` file in your project run `npx expo customize babel.config.js`.

Don't forget to clear the cache after making changes.

    npx expo start --clear

### [Authenticating Users](#authenticating-users)

With Better Auth initialized, you can now use the `authClient` to authenticate users in your Expo app.

#### [Social Sign-In](#social-sign-in)

For social sign-in, you can use the `authClient.signIn.social` method with the provider name and a callback URL.

app/social-sign-in.tsx

    import { Button } from "react-native";
    
    export default function SocialSignIn() {
        const handleLogin = async () => {
            await authClient.signIn.social({
                provider: "google",
                callbackURL: "/dashboard" // this will be converted to a deep link (eg. `myapp://dashboard`) on native
            })
        };
        return <Button title="Login with Google" onPress={handleLogin} />;
    }

#### [IdToken Sign-In](#idtoken-sign-in)

If you want to make provider request on the mobile device and then verify the ID token on the server, you can use the `authClient.signIn.social` method with the `idToken` option.

app/social-sign-in.tsx

    import { Button } from "react-native";
    
    export default function SocialSignIn() {
        const handleLogin = async () => {
            await authClient.signIn.social({
                provider: "google", // only google, apple and facebook are supported for idToken signIn
                idToken: {
                    token: "...", // ID token from provider
                    nonce: "...", // nonce from provider (optional)
                }
                callbackURL: "/dashboard" // this will be converted to a deep link (eg. `myapp://dashboard`) on native
            })
        };
        return <Button title="Login with Google" onPress={handleLogin} />;
    }

### [Session](#session)

Better Auth provides a `useSession` hook to access the current user's session in your app.

app/index.tsx

    import { Text } from "react-native";
    import { authClient } from "@/lib/auth-client";
    
    export default function Index() {
        const { data: session } = authClient.useSession();
    
        return <Text>Welcome, {session?.user.name}</Text>;
    }

On native, the session data will be cached in SecureStore. This will allow you to remove the need for a loading spinner when the app is reloaded. You can disable this behavior by passing the `disableCache` option to the client.

### [Making Authenticated Requests to Your Server](#making-authenticated-requests-to-your-server)

To make authenticated requests to your server that require the user's session, you have to retrieve the session cookie from `SecureStore` and manually add it to your request headers.

    import { authClient } from "@/lib/auth-client";
    
    const makeAuthenticatedRequest = async () => {
      const cookies = authClient.getCookie(); 
      const headers = {
        "Cookie": cookies, 
      };
      const response = await fetch("http://localhost:8081/api/secure-endpoint", { 
        headers,
        // 'include' can interfere with the cookies we just set manually in the headers
        credentials: "omit"
      });
      const data = await response.json();
      return data;
    };

**Example: Usage With TRPC**

lib/trpc-provider.tsx

    //...other imports
    import { authClient } from "@/lib/auth-client"; 
    
    export const api = createTRPCReact<AppRouter>();
    
    export function TRPCProvider(props: { children: React.ReactNode }) {
      const [queryClient] = useState(() => new QueryClient());
      const [trpcClient] = useState(() =>
        api.createClient({
          links: [
            httpBatchLink({
              //...your other options
              headers() {
                const headers = new Map<string, string>(); 
                const cookies = authClient.getCookie(); 
                if (cookies) { 
                  headers.set("Cookie", cookies); 
                } 
                return Object.fromEntries(headers); 
              },
            }),
          ],
        }),
      );
    
      return (
        <api.Provider client={trpcClient} queryClient={queryClient}>
          <QueryClientProvider client={queryClient}>
            {props.children}
          </QueryClientProvider>
        </api.Provider>
      );
    }

### [Expo Client](#expo-client)

**storage**: the storage mechanism used to cache the session data and cookies.

lib/auth-client.ts

    import { createAuthClient } from "better-auth/react";
    import SecureStorage from "expo-secure-store";
    
    const authClient = createAuthClient({
        baseURL: "http://localhost:8081",
        storage: SecureStorage
    });

**scheme**: scheme is used to deep link back to your app after a user has authenticated using oAuth providers. By default, Better Auth tries to read the scheme from the `app.json` file. If you need to override this, you can pass the scheme option to the client.

lib/auth-client.ts

    import { createAuthClient } from "better-auth/react";
    
    const authClient = createAuthClient({
        baseURL: "http://localhost:8081",
        scheme: "myapp"
    });

**disableCache**: By default, the client will cache the session data in SecureStore. You can disable this behavior by passing the `disableCache` option to the client.

lib/auth-client.ts

    import { createAuthClient } from "better-auth/react";
    
    const authClient = createAuthClient({
        baseURL: "http://localhost:8081",
        disableCache: true
    });

**cookiePrefix**: Prefix(es) for server cookie names to identify which cookies belong to better-auth. This prevents infinite refetching when third-party cookies are set. Can be a single string or an array of strings to match multiple prefixes. Defaults to `"better-auth"`.

lib/auth-client.ts

    import { createAuthClient } from "better-auth/react";
    import { expoClient } from "@better-auth/expo/client";
    import * as SecureStore from "expo-secure-store";
    
    const authClient = createAuthClient({
        baseURL: "http://localhost:8081",
        plugins: [
            expoClient({
                storage: SecureStore,
                // Single prefix
                cookiePrefix: "better-auth"
            })
        ]
    });

You can also provide multiple prefixes to match cookies from different authentication systems:

lib/auth-client.ts

    const authClient = createAuthClient({
        baseURL: "http://localhost:8081",
        plugins: [
            expoClient({
                storage: SecureStore,
                // Multiple prefixes
                cookiePrefix: ["better-auth", "my-app", "custom-auth"]
            })
        ]
    });

**Important:** If you're using plugins like passkey with a custom `webAuthnChallengeCookie` option, make sure to include the cookie prefix in the `cookiePrefix` array. For example, if you set `webAuthnChallengeCookie: "my-app-passkey"`, include `"my-app"` in your `cookiePrefix`. See the [Passkey plugin documentation](https://www.better-auth.com/docs/plugins/passkey#expo-integration) for more details.

### [Expo Servers](#expo-servers)

Server plugin options:

**disableOriginOverride**: Override the origin for Expo API routes (default: false). Enable this if you're facing cors origin issues with Expo API routes.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/expo.mdx)</content>
</page>

<page>
  <title>Generic OAuth | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/generic-oauth</url>
  <content>The Generic OAuth plugin provides a flexible way to integrate authentication with any OAuth provider. It supports both OAuth 2.0 and OpenID Connect (OIDC) flows, allowing you to easily add social login or custom OAuth authentication to your application.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

To use the Generic OAuth plugin, add it to your auth config.

auth.ts

    import { betterAuth } from "better-auth"
    import { genericOAuth } from "better-auth/plugins"
    
    export const auth = betterAuth({
        // ... other config options
        plugins: [ 
            genericOAuth({ 
                config: [ 
                    { 
                        providerId: "provider-id", 
                        clientId: "test-client-id", 
                        clientSecret: "test-client-secret", 
                        discoveryUrl: "https://auth.example.com/.well-known/openid-configuration", 
                        // ... other config options
                    }, 
                    // Add more providers as needed
                ] 
            }) 
        ]
    })

### [Add the client plugin](#add-the-client-plugin)

Include the Generic OAuth client plugin in your authentication client instance.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { genericOAuthClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            genericOAuthClient()
        ]
    })

The Generic OAuth plugin provides endpoints for initiating the OAuth flow and handling the callback. Here's how to use them:

### [Initiate OAuth Sign-In](#initiate-oauth-sign-in)

To start the OAuth sign-in process:

    const { data, error } = await authClient.signIn.oauth2({    providerId: "provider-id", // required    callbackURL: "/dashboard",    errorCallbackURL: "/error-page",    newUserCallbackURL: "/welcome",    disableRedirect: false,    scopes: ["my-scope"],    requestSignUp: false,});

| Prop | Description | Type |
| --- | --- | --- |
| `providerId` | 
The provider ID for the OAuth provider.

 | `string` |
| `callbackURL?` | 

The URL to redirect to after sign in.

 | `string` |
| `errorCallbackURL?` | 

The URL to redirect to if an error occurs.

 | `string` |
| `newUserCallbackURL?` | 

The URL to redirect to after login if the user is new.

 | `string` |
| `disableRedirect?` | 

Disable redirect.

 | `boolean` |
| `scopes?` | 

Scopes to be passed to the provider authorization request.

 | `string[]` |
| `requestSignUp?` | 

Explicitly request sign-up. Useful when disableImplicitSignUp is true for this provider.

 | `boolean` |

### [Linking OAuth Accounts](#linking-oauth-accounts)

To link an OAuth account to an existing user:

    const { data, error } = await authClient.oauth2.link({    providerId: "my-provider-id", // required    callbackURL: "/successful-link", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `providerId` | 
The OAuth provider ID.

 | `string` |
| `callbackURL` | 

The URL to redirect to once the account linking was complete.

 | `string` |

### [Handle OAuth Callback](#handle-oauth-callback)

The plugin mounts a route to handle the OAuth callback `/oauth2/callback/:providerId`. This means by default `${baseURL}/api/auth/oauth2/callback/:providerId` will be used as the callback URL. Make sure your OAuth provider is configured to use this URL.

Better Auth provides pre-configured helper functions for popular OAuth providers. These helpers handle the provider-specific configuration, including discovery URLs and user info endpoints.

### [Supported Providers](#supported-providers)

*   **Auth0** - `auth0(options)`
*   **HubSpot** - `hubspot(options)`
*   **Keycloak** - `keycloak(options)`
*   **LINE** - `line(options)`
*   **Microsoft Entra ID (Azure AD)** - `microsoftEntraId(options)`
*   **Okta** - `okta(options)`
*   **Slack** - `slack(options)`

### [Example: Using Pre-configured Providers](#example-using-pre-configured-providers)

auth.ts

    import { betterAuth } from "better-auth"
    import { genericOAuth, auth0, hubspot, keycloak, line, microsoftEntraId, okta, slack } from "better-auth/plugins"
    
    export const auth = betterAuth({
      plugins: [
        genericOAuth({
          config: [
            auth0({
              clientId: process.env.AUTH0_CLIENT_ID,
              clientSecret: process.env.AUTH0_CLIENT_SECRET,
              domain: process.env.AUTH0_DOMAIN,
            }),
            hubspot({
              clientId: process.env.HUBSPOT_CLIENT_ID,
              clientSecret: process.env.HUBSPOT_CLIENT_SECRET,
              scopes: ["oauth", "contacts"],
            }),
            keycloak({
              clientId: process.env.KEYCLOAK_CLIENT_ID,
              clientSecret: process.env.KEYCLOAK_CLIENT_SECRET,
              issuer: process.env.KEYCLOAK_ISSUER,
            }),
            // LINE supports multiple channels (countries) - use different providerIds
            line({
              providerId: "line-jp",
              clientId: process.env.LINE_JP_CLIENT_ID,
              clientSecret: process.env.LINE_JP_CLIENT_SECRET,
            }),
            line({
              providerId: "line-th",
              clientId: process.env.LINE_TH_CLIENT_ID,
              clientSecret: process.env.LINE_TH_CLIENT_SECRET,
            }),
            microsoftEntraId({
              clientId: process.env.MS_APP_ID,
              clientSecret: process.env.MS_CLIENT_SECRET,
              tenantId: process.env.MS_TENANT_ID,
            }),
            okta({
              clientId: process.env.OKTA_CLIENT_ID,
              clientSecret: process.env.OKTA_CLIENT_SECRET,
              issuer: process.env.OKTA_ISSUER,
            }),
            slack({
              clientId: process.env.SLACK_CLIENT_ID,
              clientSecret: process.env.SLACK_CLIENT_SECRET,
            }),
          ],
        }),
      ],
    })

Each provider helper accepts common OAuth options (extending `BaseOAuthProviderOptions`) plus provider-specific fields:

*   **Auth0**: Requires `domain` (e.g., `dev-xxx.eu.auth0.com`)
*   **HubSpot**: No additional required fields. Optional `scopes` (defaults to `["oauth"]`)
*   **Keycloak**: Requires `issuer` (e.g., `https://my-domain/realms/MyRealm`)
*   **LINE**: Optional `providerId` (defaults to `"line"`). LINE requires separate channels for different countries (Japan, Thailand, Taiwan, etc.), so you can call `line()` multiple times with different `providerId`s and credentials to support multiple countries
*   **Microsoft Entra ID**: Requires `tenantId` (can be a GUID, `"common"`, `"organizations"`, or `"consumers"`)
*   **Okta**: Requires `issuer` (e.g., `https://dev-xxxxx.okta.com/oauth2/default`)
*   **Slack**: No additional required fields

All providers support the same optional fields:

*   `scopes?: string[]` - Array of OAuth scopes to request
*   `redirectURI?: string` - Custom redirect URI
*   `pkce?: boolean` - Enable PKCE (defaults to `false`)
*   `disableImplicitSignUp?: boolean` - Disable automatic sign-up for new users
*   `disableSignUp?: boolean` - Disable sign-up entirely
*   `overrideUserInfo?: boolean` - Override user info on sign in

When adding the plugin to your auth config, you can configure multiple OAuth providers. You can either use the pre-configured provider helpers (shown above) or create custom configurations manually.

### [Manual Configuration](#manual-configuration)

Each provider configuration object supports the following options:

    interface GenericOAuthConfig {
      providerId: string;
      discoveryUrl?: string;
      authorizationUrl?: string;
      tokenUrl?: string;
      userInfoUrl?: string;
      clientId: string;
      clientSecret: string;
      scopes?: string[];
      redirectURI?: string;
      responseType?: string;
      prompt?: string;
      pkce?: boolean;
      accessType?: string;
      getUserInfo?: (tokens: OAuth2Tokens) => Promise<User | null>;
    }

### [Other Provider Configurations](#other-provider-configurations)

**providerId**: A unique string to identify the OAuth provider configuration.

**discoveryUrl**: (Optional) URL to fetch the provider's OAuth 2.0/OIDC configuration. If provided, endpoints like `authorizationUrl`, `tokenUrl`, and `userInfoUrl` can be auto-discovered.

**authorizationUrl**: (Optional) The OAuth provider's authorization endpoint. Not required if using `discoveryUrl`.

**tokenUrl**: (Optional) The OAuth provider's token endpoint. Not required if using `discoveryUrl`.

**userInfoUrl**: (Optional) The endpoint to fetch user profile information. Not required if using `discoveryUrl`.

**clientId**: The OAuth client ID issued by your provider.

**clientSecret**: The OAuth client secret issued by your provider.

**scopes**: (Optional) An array of scopes to request from the provider (e.g., `["openid", "email", "profile"]`).

**redirectURI**: (Optional) The redirect URI to use for the OAuth flow. If not set, a default is constructed based on your app's base URL.

**responseType**: (Optional) The OAuth response type. Defaults to `"code"` for authorization code flow.

**responseMode**: (Optional) The response mode for the authorization code request, such as `"query"` or `"form_post"`.

**prompt**: (Optional) Controls the authentication experience (e.g., force login, consent, etc.).

**pkce**: (Optional) If true, enables PKCE (Proof Key for Code Exchange) for enhanced security. Defaults to `false`.

**accessType**: (Optional) The access type for the authorization request. Use `"offline"` to request a refresh token.

**getToken**: (Optional) A custom function to exchange authorization code for tokens. If provided, this function will be used instead of the default token exchange logic. This is useful for providers with non-standard token endpoints that use GET requests or custom parameters.

**getUserInfo**: (Optional) A custom function to fetch user info from the provider, given the OAuth tokens. If not provided, a default fetch is used.

**mapProfileToUser**: (Optional) A function to map the provider's user profile to your app's user object. Useful for custom field mapping or transformations.

**authorizationUrlParams**: (Optional) Additional query parameters to add to the authorization URL. These can override default parameters. You can also provide a function that returns the parameters.

**tokenUrlParams**: (Optional) Additional query parameters to add to the token URL. These can override default parameters. You can also provide a function that returns the parameters.

**disableImplicitSignUp**: (Optional) If true, disables automatic sign-up for new users. Sign-in must be explicitly requested with sign-up intent.

**disableSignUp**: (Optional) If true, disables sign-up for new users entirely. Only existing users can sign in.

**authentication**: (Optional) The authentication method for token requests. Can be `'basic'` or `'post'`. Defaults to `'post'`.

**discoveryHeaders**: (Optional) Custom headers to include in the discovery request. Useful for providers that require special headers.

**authorizationHeaders**: (Optional) Custom headers to include in the authorization request. Useful for providers that require special headers.

**overrideUserInfo**: (Optional) If true, the user's info in your database will be updated with the provider's info every time they sign in. Defaults to `false`.

### [Custom Token Exchange](#custom-token-exchange)

For providers with non-standard token endpoints that use GET requests or custom parameters, you can provide a custom `getToken` function:

    genericOAuth({
      config: [
        {
          providerId: "custom-provider",
          clientId: process.env.CUSTOM_CLIENT_ID!,
          clientSecret: process.env.CUSTOM_CLIENT_SECRET,
          authorizationUrl: "https://provider.example.com/oauth/authorize",
          scopes: ["profile", "email"],
          // Custom token exchange for non-standard endpoints
          getToken: async ({ code, redirectURI }) => {
            // Example: GET request instead of POST
            const response = await fetch(
              `https://provider.example.com/oauth/token?` +
              `client_id=${process.env.CUSTOM_CLIENT_ID}&` +
              `client_secret=${process.env.CUSTOM_CLIENT_SECRET}&` +
              `code=${code}&` +
              `redirect_uri=${redirectURI}&` +
              `grant_type=authorization_code`,
              { method: "GET" }
            );
    
            const data = await response.json();
    
            return {
              accessToken: data.access_token,
              refreshToken: data.refresh_token,
              accessTokenExpiresAt: new Date(Date.now() + data.expires_in * 1000),
              scopes: data.scope?.split(" ") ?? [],
              // Preserve provider-specific fields in raw
              raw: data,
            };
          },
          getUserInfo: async (tokens) => {
            // Access provider-specific fields from raw token data
            const userId = tokens.raw?.user_id as string;
    
            const response = await fetch(
              `https://provider.example.com/api/user?` +
              `access_token=${tokens.accessToken}`
            );
    
            const data = await response.json();
    
            return {
              id: userId,
              name: data.display_name,
              email: data.email,
              image: data.avatar_url,
              emailVerified: data.email_verified,
            };
          },
        },
      ],
    });

### [Custom User Info Fetching](#custom-user-info-fetching)

You can provide a custom `getUserInfo` function to handle specific provider requirements:

    genericOAuth({
      config: [
        {
          providerId: "custom-provider",
          // ... other config options
          getUserInfo: async (tokens) => {
            // Custom logic to fetch and return user info
            const userInfo = await fetchUserInfoFromCustomProvider(tokens);
            return {
              id: userInfo.sub,
              email: userInfo.email,
              name: userInfo.name,
              // ... map other fields as needed
            };
          }
        }
      ]
    })

### [Map User Info Fields](#map-user-info-fields)

If the user info returned by the provider does not match the expected format, or you need to map additional fields, you can use the `mapProfileToUser`:

    genericOAuth({
      config: [
        {
          providerId: "custom-provider",
          // ... other config options
          mapProfileToUser: async (profile) => {
            return {
              firstName: profile.given_name,
              // ... map other fields as needed
            };
          }
        }
      ]
    })

### [Accessing Raw Token Data](#accessing-raw-token-data)

The `tokens` parameter includes a `raw` field that preserves the original token response from the provider. This is useful for accessing provider-specific fields:

    getUserInfo: async (tokens) => {
      // Access provider-specific fields
      const customField = tokens.raw?.custom_provider_field as string;
      const userId = tokens.raw?.provider_user_id as string;
    
      // Use in your logic
      return {
        id: userId,
        // ...
      };
    }

### [Error Handling](#error-handling)

The plugin includes built-in error handling for common OAuth issues. Errors are typically redirected to your application's error page with an appropriate error message in the URL parameters. If the callback URL is not provided, the user will be redirected to Better Auth's default error page.</content>
</page>

<page>
  <title>Anonymous | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/anonymous</url>
  <content>The Anonymous plugin allows users to have an authenticated experience without requiring them to provide an email address, password, OAuth provider, or any other Personally Identifiable Information (PII). Users can later link an authentication method to their account when ready.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

To enable anonymous authentication, add the anonymous plugin to your authentication configuration.

auth.ts

    import { betterAuth } from "better-auth"
    import { anonymous } from "better-auth/plugins"
    
    export const auth = betterAuth({
        // ... other config options
        plugins: [
            anonymous() 
        ]
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

Next, include the anonymous client plugin in your authentication client instance.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { anonymousClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            anonymousClient()
        ]
    })

### [Sign In](#sign-in)

To sign in a user anonymously, use the `signIn.anonymous()` method.

example.ts

    const user = await authClient.signIn.anonymous()

### [Link Account](#link-account)

If a user is already signed in anonymously and tries to `signIn` or `signUp` with another method, their anonymous activities can be linked to the new account.

To do that you first need to provide `onLinkAccount` callback to the plugin.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        plugins: [
            anonymous({
                onLinkAccount: async ({ anonymousUser, newUser }) => {
                   // perform actions like moving the cart items from anonymous user to the new user
                }
            })
        ]

Then when you call `signIn` or `signUp` with another method, the `onLinkAccount` callback will be called. And the `anonymousUser` will be deleted by default.

example.ts

    const user = await authClient.signIn.email({
        email,
    })

### [`emailDomainName`](#emaildomainname)

The domain name to use when generating an email address for anonymous users. If not provided, the default format `temp@{id}.com` is used.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        plugins: [
            anonymous({
                emailDomainName: "example.com" // -> temp-{id}@example.com
            })
        ]
    })

### [`generateRandomEmail`](#generaterandomemail)

A custom function to generate email addresses for anonymous users. This allows you to define your own email format. The function can be synchronous or asynchronous.

auth.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        plugins: [
            anonymous({
                generateRandomEmail: () => { 
                    const id = crypto.randomUUID() 
                    return `guest-${id}@example.com`
                } 
            })
        ]
    })

**Notes:**

*   If `generateRandomEmail` is provided, `emailDomainName` is ignored.
*   You are responsible for ensuring the email is unique to avoid conflicts. The returned email must be in a valid format.

### [`onLinkAccount`](#onlinkaccount)

A callback function that is called when an anonymous user links their account to a new authentication method. The callback receives an object with the `anonymousUser` and the `newUser`.

### [`disableDeleteAnonymousUser`](#disabledeleteanonymoususer)

By default, the anonymous user is deleted when the account is linked to a new authentication method. Set this option to `true` to disable this behavior.

### [`generateName`](#generatename)

A callback function that is called to generate a name for the anonymous user. Useful if you want to have random names for anonymous users, or if `name` is unique in your database.

The anonymous plugin requires an additional field in the user table:

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| isAnonymous | boolean |  | Indicates whether the user is anonymous. |

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/anonymous.mdx)</content>
</page>

<page>
  <title>Phone Number | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/phone-number</url>
  <content>The phone number plugin extends the authentication system by allowing users to sign in and sign up using their phone number. It includes OTP (One-Time Password) functionality to verify phone numbers.

### [Add Plugin to the server](#add-plugin-to-the-server)

auth.ts

    import { betterAuth } from "better-auth"
    import { phoneNumber } from "better-auth/plugins"
    
    const auth = betterAuth({
        plugins: [ 
            phoneNumber({  
                sendOTP: ({ phoneNumber, code }, ctx) => { 
                    // Implement sending OTP code via SMS
                } 
            }) 
        ] 
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { phoneNumberClient } from "better-auth/client/plugins"
    
    const authClient = createAuthClient({
        plugins: [ 
            phoneNumberClient() 
        ] 
    })

### [Send OTP for Verification](#send-otp-for-verification)

To send an OTP to a user's phone number for verification, you can use the `sendVerificationCode` endpoint.

POST

/phone-number/send-otp

    const { data, error } = await authClient.phoneNumber.sendOtp({    phoneNumber: "+1234567890", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `phoneNumber` | 
Phone number to send OTP.

 | `string` |

### [Verify Phone Number](#verify-phone-number)

After the OTP is sent, users can verify their phone number by providing the code.

    const { data, error } = await authClient.phoneNumber.verify({    phoneNumber: "+1234567890", // required    code: "123456", // required    disableSession: false,    updatePhoneNumber: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `phoneNumber` | 
Phone number to verify.

 | `string` |
| `code` | 

OTP code.

 | `string` |
| `disableSession?` | 

Disable session creation after verification.

 | `boolean` |
| `updatePhoneNumber?` | 

Check if there is a session and update the phone number.

 | `boolean` |

When the phone number is verified, the `phoneNumberVerified` field in the user table is set to `true`. If `disableSession` is not set to `true`, a session is created for the user. Additionally, if `callbackOnVerification` is provided, it will be called.

### [Allow Sign-Up with Phone Number](#allow-sign-up-with-phone-number)

To allow users to sign up using their phone number, you can pass `signUpOnVerification` option to your plugin configuration. It requires you to pass `getTempEmail` function to generate a temporary email for the user.

auth.ts

    export const auth = betterAuth({
        plugins: [
            phoneNumber({
                sendOTP: ({ phoneNumber, code }, ctx) => {
                    // Implement sending OTP code via SMS
                },
                signUpOnVerification: {
                    getTempEmail: (phoneNumber) => {
                        return `${phoneNumber}@my-site.com`
                    },
                    //optionally, you can also pass `getTempName` function to generate a temporary name for the user
                    getTempName: (phoneNumber) => {
                        return phoneNumber //by default, it will use the phone number as the name
                    }
                }
            })
        ]
    })

### [Sign In with Phone Number](#sign-in-with-phone-number)

In addition to signing in a user using send-verify flow, you can also use phone number as an identifier and sign in a user using phone number and password.

POST

/sign-in/phone-number

    const { data, error } = await authClient.signIn.phoneNumber({    phoneNumber: "+1234567890", // required    password, // required    rememberMe: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `phoneNumber` | 
Phone number to sign in.

 | `string` |
| `password` | 

Password to use for sign in.

 | `string` |
| `rememberMe?` | 

Remember the session.

 | `boolean` |

### [Update Phone Number](#update-phone-number)

Updating phone number uses the same process as verifying a phone number. The user will receive an OTP code to verify the new phone number.

auth-client.ts

    await authClient.phoneNumber.sendOtp({
        phoneNumber: "+1234567890" // New phone number
    })

Then verify the new phone number with the OTP code.

auth-client.ts

    const isVerified = await authClient.phoneNumber.verify({
        phoneNumber: "+1234567890",
        code: "123456",
        updatePhoneNumber: true // Set to true to update the phone number
    })

If a user session exist the phone number will be updated automatically.

### [Disable Session Creation](#disable-session-creation)

By default, the plugin creates a session for the user after verifying the phone number. You can disable this behavior by passing `disableSession: true` to the `verify` method.

auth-client.ts

    const isVerified = await authClient.phoneNumber.verify({
        phoneNumber: "+1234567890",
        code: "123456",
        disableSession: true
    })

### [Request Password Reset](#request-password-reset)

To initiate a request password reset flow using `phoneNumber`, you can start by calling `requestPasswordReset` on the client to send an OTP code to the user's phone number.

POST

/phone-number/request-password-reset

    const { data, error } = await authClient.phoneNumber.requestPasswordReset({    phoneNumber: "+1234567890", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `phoneNumber` | 
The phone number which is associated with the user.

 | `string` |

Then, you can reset the password by calling `resetPassword` on the client with the OTP code and the new password.

POST

/phone-number/reset-password

    const { data, error } = await authClient.phoneNumber.resetPassword({    otp: "123456", // required    phoneNumber: "+1234567890", // required    newPassword: "new-and-secure-password", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `otp` | 
The one time password to reset the password.

 | `string` |
| `phoneNumber` | 

The phone number to the account which intends to reset the password for.

 | `string` |
| `newPassword` | 

The new password.

 | `string` |

### [`otpLength`](#otplength)

The length of the OTP code to be generated. Default is `6`.

### [`sendOTP`](#sendotp)

A function that sends the OTP code to the user's phone number. It takes the phone number and the OTP code as arguments.

### [`expiresIn`](#expiresin)

The time in seconds after which the OTP code expires. Default is `300` seconds.

### [`callbackOnVerification`](#callbackonverification)

A function that is called after the phone number is verified. It takes the phone number and the user object as the first argument and a request object as the second argument.

    export const auth = betterAuth({
        plugins: [
            phoneNumber({
                sendOTP: ({ phoneNumber, code }, ctx) => {
                    // Implement sending OTP code via SMS
                },
                callbackOnVerification: async ({ phoneNumber, user }, ctx) => { 
                    // Implement callback after phone number verification
                } 
            })
        ]
    })

### [`sendPasswordResetOTP`](#sendpasswordresetotp)

A function that sends the OTP code to the user's phone number for password reset. It takes the phone number and the OTP code as arguments.

### [`phoneNumberValidator`](#phonenumbervalidator)

A custom function to validate the phone number. It takes the phone number as an argument and returns a boolean indicating whether the phone number is valid.

### [`verifyOTP`](#verifyotp)

A custom function to verify the OTP code. When provided, this function will be used instead of the default internal verification logic. This is useful when you want to integrate with external SMS providers that handle OTP verification (e.g., Twilio Verify, AWS SNS). The function takes an object with `phoneNumber` and `code` properties and a request object, and returns a boolean or a promise that resolves to a boolean indicating whether the OTP is valid.

    export const auth = betterAuth({
        plugins: [
            phoneNumber({
                sendOTP: ({ phoneNumber, code }, ctx) => {
                    // Send OTP via your SMS provider
                },
                verifyOTP: async ({ phoneNumber, code }, ctx) => { 
                    // Verify OTP with your desired logic (e.g., Twilio Verify)
                    // This is just an example, not a real implementation.
                    const isValid = await twilioClient.verify 
                        .services('YOUR_SERVICE_SID') 
                        .verificationChecks 
                        .create({ to: phoneNumber, code }); 
                    return isValid.status === 'approved'; 
                } 
            })
        ]
    })

When using this option, ensure that proper validation is implemented, as it overrides our internal verification logic.

### [`signUpOnVerification`](#signuponverification)

An object with the following properties:

*   `getTempEmail`: A function that generates a temporary email for the user. It takes the phone number as an argument and returns the temporary email.
*   `getTempName`: A function that generates a temporary name for the user. It takes the phone number as an argument and returns the temporary name.

### [`requireVerification`](#requireverification)

When enabled, users cannot sign in with their phone number until it has been verified. If an unverified user attempts to sign in, the server will respond with a 401 error (PHONE\_NUMBER\_NOT\_VERIFIED) and automatically trigger an OTP send to start the verification process.

The plugin requires 2 fields to be added to the user table

### [User Table](#user-table)

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| phoneNumber | string |  | The phone number of the user |
| phoneNumberVerified | boolean |  | Whether the phone number is verified or not |

### [OTP Verification Attempts](#otp-verification-attempts)

The phone number plugin includes a built-in protection against brute force attacks by limiting the number of verification attempts for each OTP code.

    phoneNumber({
      allowedAttempts: 3, // default is 3
      // ... other options
    })

When a user exceeds the allowed number of verification attempts:

*   The OTP code is automatically deleted
*   Further verification attempts will return a 403 (Forbidden) status with "Too many attempts" message
*   The user will need to request a new OTP code to continue

Example error response after exceeding attempts:

    {
      "error": {
        "status": 403,
        "message": "Too many attempts"
      }
    }

When receiving a 403 status, prompt the user to request a new OTP code</content>
</page>

<page>
  <title>Hono Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/hono</url>
  <content>Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Mount the handler](#mount-the-handler)

We need to mount the handler to Hono endpoint.

    import { Hono } from "hono";
    import { auth } from "./auth";
    import { serve } from "@hono/node-server";
    
    const app = new Hono();
    
    app.on(["POST", "GET"], "/api/auth/*", (c) => {
    	return auth.handler(c.req.raw);
    });
    
    serve(app);

### [Cors](#cors)

To configure cors, you need to use the `cors` plugin from `hono/cors`.

    import { Hono } from "hono";
    import { auth } from "./auth";
    import { serve } from "@hono/node-server";
    import { cors } from "hono/cors";
     
    const app = new Hono();
    
    app.use(
    	"/api/auth/*", // or replace with "*" to enable cors for all routes
    	cors({
    		origin: "http://localhost:3001", // replace with your origin
    		allowHeaders: ["Content-Type", "Authorization"],
    		allowMethods: ["POST", "GET", "OPTIONS"],
    		exposeHeaders: ["Content-Length"],
    		maxAge: 600,
    		credentials: true,
    	}),
    );
    
    app.on(["POST", "GET"], "/api/auth/*", (c) => {
    	return auth.handler(c.req.raw);
    });
    
    serve(app);

> **Important:** CORS middleware must be registered before your routes. This ensures that cross-origin requests are properly handled before they reach your authentication endpoints.

### [Middleware](#middleware)

You can add a middleware to save the `session` and `user` in a `context` and also add validations for every route.

    import { Hono } from "hono";
    import { auth } from "./auth";
    import { serve } from "@hono/node-server";
    import { cors } from "hono/cors";
     
    const app = new Hono<{
    	Variables: {
    		user: typeof auth.$Infer.Session.user | null;
    		session: typeof auth.$Infer.Session.session | null
    	}
    }>();
    
    app.use("*", async (c, next) => {
    	const session = await auth.api.getSession({ headers: c.req.raw.headers });
    
      	if (!session) {
        	c.set("user", null);
        	c.set("session", null);
        	await next();
            return;
      	}
    
      	c.set("user", session.user);
      	c.set("session", session.session);
      	await next();
    });
    
    app.on(["POST", "GET"], "/api/auth/*", (c) => {
    	return auth.handler(c.req.raw);
    });
    
    
    serve(app);

This will allow you to access the `user` and `session` object in all of your routes.

    app.get("/session", (c) => {
    	const session = c.get("session")
    	const user = c.get("user")
    	
    	if(!user) return c.body(null, 401);
    
      	return c.json({
    	  session,
    	  user
    	});
    });

### [Cross-Domain Cookies](#cross-domain-cookies)

By default, all Better Auth cookies are set with `SameSite=Lax`. If you need to use cookies across different domains, you’ll need to set `SameSite=None` and `Secure=true`. However, we recommend using subdomains whenever possible, as this allows you to keep `SameSite=Lax`. To enable cross-subdomain cookies, simply turn on `crossSubDomainCookies` in your auth config.

auth.ts

    export const auth = createAuth({
      advanced: {
        crossSubDomainCookies: {
          enabled: true
        }
      }
    })

If you still need to set `SameSite=None` and `Secure=true`, you can adjust these attributes globally through `cookieOptions` in the `createAuth` configuration.

auth.ts

    export const auth = createAuth({
      advanced: {
        defaultCookieAttributes: {
          sameSite: "none",
          secure: true,
          partitioned: true // New browser standards will mandate this for foreign cookies
        }
      }
    })

You can also customize cookie attributes individually by setting them within `cookies` in your auth config.

auth.ts

    export const auth = createAuth({
      advanced: {
        cookies: {
          sessionToken: {
            attributes: {
              sameSite: "none",
              secure: true,
              partitioned: true // New browser standards will mandate this for foreign cookies
            }
          }
        }
      }
    })

### [Client-Side Configuration](#client-side-configuration)

When using the Hono client (`@hono/client`) to make requests to your Better Auth-protected endpoints, you need to configure it to send credentials (cookies) with cross-origin requests.

api.ts

    import { hc } from "hono/client";
    import type { AppType } from "./server"; // Your Hono app type
    
    const client = hc<AppType>("http://localhost:8787/", {
      init: {
        credentials: "include", // Required for sending cookies cross-origin
      },
    });
    
    // Now your client requests will include credentials
    const response = await client.someProtectedEndpoint.$get();

This configuration is necessary when:

*   Your client and server are on different domains/ports during development
*   You're making cross-origin requests in production
*   You need to send authentication cookies with your requests

The `credentials: "include"` option tells the fetch client to send cookies even for cross-origin requests. This works in conjunction with the CORS configuration on your server that has `credentials: true`.

> **Note:** Make sure your CORS configuration on the server matches your client's domain, and that `credentials: true` is set in both the server's CORS config and the client's fetch config.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/hono.mdx)</content>
</page>

<page>
  <title>Express Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/express</url>
  <content>This guide will show you how to integrate Better Auth with [express.js](https://expressjs.com/).

Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

Note that CommonJS (cjs) isn't supported. Use ECMAScript Modules (ESM) by setting `"type": "module"` in your `package.json` or configuring your `tsconfig.json` to use ES modules.

### [Mount the handler](#mount-the-handler)

To enable Better Auth to handle requests, we need to mount the handler to an API route. Create a catch-all route to manage all requests to `/api/auth/*` in case of ExpressJS v4 or `/api/auth/*splat` in case of ExpressJS v5 (or any other path specified in your Better Auth options).

Don’t use `express.json()` before the Better Auth handler. Use it only for other routes, or the client API will get stuck on "pending".

server.ts

    import express from "express";
    import { toNodeHandler } from "better-auth/node";
    import { auth } from "./auth";
    
    const app = express();
    const port = 3005;
    
    app.all("/api/auth/*", toNodeHandler(auth)); // For ExpressJS v4
    // app.all("/api/auth/*splat", toNodeHandler(auth)); For ExpressJS v5 
    
    // Mount express json middleware after Better Auth handler
    // or only apply it to routes that don't interact with Better Auth
    app.use(express.json());
    
    app.listen(port, () => {
    	console.log(`Example app listening on port ${port}`);
    });

After completing the setup, start your server. Better Auth will be ready to use. You can send a `GET` request to the `/ok` endpoint (`/api/auth/ok`) to verify that the server is running.

### [Cors Configuration](#cors-configuration)

To add CORS (Cross-Origin Resource Sharing) support to your Express server when integrating Better Auth, you can use the `cors` middleware. Below is an updated example showing how to configure CORS for your server:

    import express from "express";
    import cors from "cors"; // Import the CORS middleware
    import { toNodeHandler, fromNodeHeaders } from "better-auth/node";
    import { auth } from "./auth";
    
    const app = express();
    const port = 3005;
    
    // Configure CORS middleware
    app.use(
      cors({
        origin: "http://your-frontend-domain.com", // Replace with your frontend's origin
        methods: ["GET", "POST", "PUT", "DELETE"], // Specify allowed HTTP methods
        credentials: true, // Allow credentials (cookies, authorization headers, etc.)
      })
    );

### [Getting the User Session](#getting-the-user-session)

To retrieve the user's session, you can use the `getSession` method provided by the `auth` object. This method requires the request headers to be passed in a specific format. To simplify this process, Better Auth provides a `fromNodeHeaders` helper function that converts Node.js request headers to the format expected by Better Auth (a `Headers` object).

Here's an example of how to use `getSession` in an Express route:

server.ts

    import { fromNodeHeaders } from "better-auth/node";
    import { auth } from "./auth"; // Your Better Auth instance
    
    app.get("/api/me", async (req, res) => {
     	const session = await auth.api.getSession({
          headers: fromNodeHeaders(req.headers),
        });
    	return res.json(session);
    });</content>
</page>

<page>
  <title>Nuxt Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/nuxt</url>
  <content>Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Create API Route](#create-api-route)

We need to mount the handler to an API route. Create a file inside `/server/api/auth` called `[...all].ts` and add the following code:

server/api/auth/\[...all\].ts

    import { auth } from "~/lib/auth"; // import your auth config
    
    export default defineEventHandler((event) => {
    	return auth.handler(toWebRequest(event));
    });

You can change the path on your better-auth configuration but it's recommended to keep it as `/api/auth/[...all]`

### [Migrate the database](#migrate-the-database)

Run the following command to create the necessary tables in your database:

    npx @better-auth/cli migrate

Create a client instance. You can name the file anything you want. Here we are creating `client.ts` file inside the `lib/` directory.

auth-client.ts

    import { createAuthClient } from "better-auth/vue" // make sure to import from better-auth/vue
    
    export const authClient = createAuthClient({
        //you can pass client configuration here
    })

Once you have created the client, you can use it to sign up, sign in, and perform other actions. Some of the actions are reactive.

### [Example usage](#example-usage)

index.vue

    <script setup lang="ts">
    import { authClient } from "~/lib/client"
    const session = authClient.useSession()
    </script>
    
    <template>
        <div>
            <button v-if="!session?.data" @click="() => authClient.signIn.social({
                provider: 'github'
            })">
                Continue with GitHub
            </button>
            <div>
                <pre>{{ session.data }}</pre>
                <button v-if="session.data" @click="authClient.signOut()">
                    Sign out
                </button>
            </div>
        </div>
    </template>

### [Server Usage](#server-usage)

The `api` object exported from the auth instance contains all the actions that you can perform on the server. Every endpoint made inside Better Auth is a invocable as a function. Including plugins endpoints.

**Example: Getting Session on a server API route**

server/api/example.ts

    import { auth } from "~/lib/auth";
    
    export default defineEventHandler((event) => {
        const session = await auth.api.getSession({
          headers: event.headers
        });
    
       if(session) {
         // access the session.session && session.user
       }
    });

### [SSR Usage](#ssr-usage)

If you are using Nuxt with SSR, you can use the `useSession` function in the `setup` function of your page component and pass `useFetch` to make it work with SSR.

index.vue

    <script setup lang="ts">
    import { authClient } from "~/lib/auth-client";
    
    const { data: session } = await authClient.useSession(useFetch);
    </script>
    
    <template>
        <p>
            {{ session }}
        </p>
    </template>

### [Middleware](#middleware)

To add middleware to your Nuxt project, you can use the `useSession` method from the client.

middleware/auth.global.ts

    import { authClient } from "~/lib/auth-client";
    export default defineNuxtRouteMiddleware(async (to, from) => {
    	const { data: session } = await authClient.useSession(useFetch); 
    	if (!session.value) {
    		if (to.path === "/dashboard") {
    			return navigateTo("/");
    		}
    	}
    });

### [Resources & Examples](#resources--examples)

*   [Nuxt and Nuxt Hub example](https://github.com/atinux/nuxthub-better-auth) on GitHub.
*   [NuxtZzle is Nuxt,Drizzle ORM example](https://github.com/leamsigc/nuxt-better-auth-drizzle) on GitHub [preview](https://nuxt-better-auth.giessen.dev/)
*   [Nuxt example](https://stackblitz.com/github/better-auth/examples/tree/main/nuxt-example) on StackBlitz.
*   [NuxSaaS (Github)](https://github.com/NuxSaaS/NuxSaaS) is a full-stack SaaS Starter Kit that leverages Better Auth for secure and efficient user authentication. [Demo](https://nuxsaas.com/)
*   [NuxtOne (Github)](https://github.com/nuxtone/nuxt-one) is a Nuxt-based starter template for building AIaaS (AI-as-a-Service) applications [preview](https://www.one.devv.zone/)

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/nuxt.mdx)</content>
</page>

<page>
  <title>Last Login Method | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/last-login-method</url>
  <content>The last login method plugin tracks the most recent authentication method used by users (email, OAuth providers, etc.). This enables you to display helpful indicators on login pages, such as "Last signed in with Google" or prioritize certain login methods based on user preferences.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { lastLoginMethod } from "better-auth/plugins"
    
    export const auth = betterAuth({
        // ... other config options
        plugins: [
            lastLoginMethod() 
        ]
    })

### [Add the client plugin to your auth client](#add-the-client-plugin-to-your-auth-client)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { lastLoginMethodClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            lastLoginMethodClient() 
        ]
    })

Once installed, the plugin automatically tracks the last authentication method used by users. You can then retrieve and display this information in your application.

### [Getting the Last Used Method](#getting-the-last-used-method)

The client plugin provides several methods to work with the last login method:

app.tsx

    import { authClient } from "@/lib/auth-client"
    
    // Get the last used login method
    const lastMethod = authClient.getLastUsedLoginMethod()
    console.log(lastMethod) // "google", "email", "github", etc.
    
    // Check if a specific method was last used
    const wasGoogle = authClient.isLastUsedLoginMethod("google")
    
    // Clear the stored method
    authClient.clearLastUsedLoginMethod()

### [UI Integration Example](#ui-integration-example)

Here's how to use the plugin to enhance your login page:

sign-in.tsx

    import { authClient } from "@/lib/auth-client"
    import { Button } from "@/components/ui/button"
    import { Badge } from "@/components/ui/badge"
    
    export function SignInPage() {
        const lastMethod = authClient.getLastUsedLoginMethod()
        
        return (
            <div className="space-y-4">
                <h1>Sign In</h1>
                
                {/* Email sign in */}
                <div className="relative">
                    <Button 
                        onClick={() => authClient.signIn.email({...})}
                        variant={lastMethod === "email" ? "default" : "outline"}
                        className="w-full"
                    >
                        Sign in with Email
                        {lastMethod === "email" && (
                            <Badge className="ml-2">Last used</Badge>
                        )}
                    </Button>
                </div>
                
                {/* OAuth providers */}
                <div className="relative">
                    <Button 
                        onClick={() => authClient.signIn.social({ provider: "google" })}
                        variant={lastMethod === "google" ? "default" : "outline"}
                        className="w-full"
                    >
                        Continue with Google
                        {lastMethod === "google" && (
                            <Badge className="ml-2">Last used</Badge>
                        )}
                    </Button>
                </div>
                
                <div className="relative">
                    <Button 
                        onClick={() => authClient.signIn.social({ provider: "github" })}
                        variant={lastMethod === "github" ? "default" : "outline"}
                        className="w-full"
                    >
                        Continue with GitHub
                        {lastMethod === "github" && (
                            <Badge className="ml-2">Last used</Badge>
                        )}
                    </Button>
                </div>
            </div>
        )
    }

By default, the last login method is stored only in cookies. For more persistent tracking and analytics, you can enable database storage.

### [Enable database storage](#enable-database-storage)

Set `storeInDatabase` to `true` in your plugin configuration:

auth.ts

    import { betterAuth } from "better-auth"
    import { lastLoginMethod } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [
            lastLoginMethod({
                storeInDatabase: true
            })
        ]
    })

### [Run database migration](#run-database-migration)

The plugin will automatically add a `lastLoginMethod` field to your user table. Run the migration to apply the changes:

### [Access database field](#access-database-field)

When database storage is enabled, the `lastLoginMethod` field becomes available in user objects:

user-profile.tsx

    import { auth } from "@/lib/auth"
    
    // Server-side access
    const session = await auth.api.getSession({ headers })
    console.log(session?.user.lastLoginMethod) // "google", "email", etc.
    
    // Client-side access via session
    const { data: session } = authClient.useSession()
    console.log(session?.user.lastLoginMethod)

### [Database Schema](#database-schema)

When `storeInDatabase` is enabled, the plugin adds the following field to the `user` table:

Table: `user`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| lastLoginMethod | string |  | The last authentication method used by the user |

### [Custom Schema Configuration](#custom-schema-configuration)

You can customize the database field name:

auth.ts

    import { betterAuth } from "better-auth"
    import { lastLoginMethod } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [
            lastLoginMethod({
                storeInDatabase: true,
                schema: {
                    user: {
                        lastLoginMethod: "last_auth_method" // Custom field name
                    }
                }
            })
        ]
    })

The last login method plugin accepts the following options:

### [Server Options](#server-options)

auth.ts

    import { betterAuth } from "better-auth"
    import { lastLoginMethod } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [
            lastLoginMethod({
                // Cookie configuration
                cookieName: "better-auth.last_used_login_method", // Default: "better-auth.last_used_login_method"
                maxAge: 60 * 60 * 24 * 30, // Default: 30 days in seconds
                
                // Database persistence
                storeInDatabase: false, // Default: false
                
                // Custom method resolution
                customResolveMethod: (ctx) => {
                    // Custom logic to determine the login method
                    if (ctx.path === "/oauth/callback/custom-provider") {
                        return "custom-provider"
                    }
                    // Return null to use default resolution
                    return null
                },
                
                // Schema customization (when storeInDatabase is true)
                schema: {
                    user: {
                        lastLoginMethod: "custom_field_name"
                    }
                }
            })
        ]
    })

**cookieName**: `string`

*   The name of the cookie used to store the last login method
*   Default: `"better-auth.last_used_login_method"`
*   **Note**: This cookie is `httpOnly: false` to allow client-side JavaScript access for UI features

**maxAge**: `number`

*   Cookie expiration time in seconds
*   Default: `2592000` (30 days)

**storeInDatabase**: `boolean`

*   Whether to store the last login method in the database
*   Default: `false`
*   When enabled, adds a `lastLoginMethod` field to the user table

**customResolveMethod**: `(ctx: GenericEndpointContext) => string | null`

*   Custom function to determine the login method from the request context
*   Return `null` to use the default resolution logic
*   Useful for custom OAuth providers or authentication flows

**schema**: `object`

*   Customize database field names when `storeInDatabase` is enabled
*   Allows mapping the `lastLoginMethod` field to a custom column name

### [Client Options](#client-options)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { lastLoginMethodClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            lastLoginMethodClient({
                cookieName: "better-auth.last_used_login_method" // Default: "better-auth.last_used_login_method"
            })
        ]
    })

**cookieName**: `string`

*   The name of the cookie to read the last login method from
*   Must match the server-side `cookieName` configuration
*   Default: `"better-auth.last_used_login_method"`

### [Default Method Resolution](#default-method-resolution)

By default, the plugin tracks these authentication methods:

*   **Email authentication**: `"email"`
*   **OAuth providers**: Provider ID (e.g., `"google"`, `"github"`, `"discord"`)
*   **OAuth2 callbacks**: Provider ID from URL path
*   **Sign up methods**: Tracked the same as sign in methods

The plugin automatically detects the method from these endpoints:

*   `/callback/:id` - OAuth callback with provider ID
*   `/oauth2/callback/:id` - OAuth2 callback with provider ID
*   `/sign-in/email` - Email sign in
*   `/sign-up/email` - Email sign up

[Cross-Domain Support](#cross-domain-support)
---------------------------------------------

The plugin automatically inherits cookie settings from Better Auth's centralized cookie system. This solves the problem where the last login method wouldn't persist across:

*   **Cross-subdomain setups**: `auth.example.com` → `app.example.com`
*   **Cross-origin setups**: `api.company.com` → `app.different.com`

When you enable `crossSubDomainCookies` or `crossOriginCookies` in your Better Auth config, the plugin will automatically use the same domain, secure, and sameSite settings as your session cookies, ensuring consistent behavior across your application.

### [Custom Provider Tracking](#custom-provider-tracking)

If you have custom OAuth providers or authentication methods, you can use the `customResolveMethod` option:

auth.ts

    import { betterAuth } from "better-auth"
    import { lastLoginMethod } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [
            lastLoginMethod({
                customResolveMethod: (ctx) => {
                    // Track custom SAML provider
                    if (ctx.path === "/saml/callback") {
                        return "saml"
                    }
                    
                    // Track magic link authentication
                    if (ctx.path === "/magic-link/verify") {
                        return "magic-link"
                    }
                    
                    // Track phone authentication
                    if (ctx.path === "/sign-in/phone") {
                        return "phone"
                    }
                    
                    // Return null to use default logic
                    return null
                }
            })
        ]
    })

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/last-login-method.mdx)</content>
</page>

<page>
  <title>OAuth Proxy | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/oauth-proxy</url>
  <content>A proxy plugin, that allows you to proxy OAuth requests. Useful for development and preview deployments where the redirect URL can't be known in advance to add to the OAuth provider.

### [Add the plugin to your **auth** config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { oAuthProxy } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [ 
            oAuthProxy({ 
                productionURL: "https://my-main-app.com", // Optional - if the URL isn't inferred correctly
                currentURL: "http://localhost:3000", // Optional - if the URL isn't inferred correctly
            }), 
        ] 
    })

### [Add redirect URL to your OAuth provider](#add-redirect-url-to-your-oauth-provider)

For the proxy server to work properly, you’ll need to pass the redirect URL of your main production app registered with the OAuth provider in your social provider config. This needs to be done for each social provider you want to proxy requests for.

    export const auth = betterAuth({
       plugins: [
           oAuthProxy(),
       ], 
       socialProviders: {
            github: {
                clientId: "your-client-id",
                clientSecret: "your-client-secret",
                redirectURI: "https://my-main-app.com/api/auth/callback/github"
            }
       }
    })

The plugin adds an endpoint to your server that proxies OAuth requests. When you initiate a social sign-in, it sets the redirect URL to this proxy endpoint. After the OAuth provider redirects back to your server, the plugin then forwards the user to the original callback URL.

    await authClient.signIn.social({
        provider: "github",
        callbackURL: "/dashboard" // the plugin will override this to something like "http://localhost:3000/api/auth/oauth-proxy?callbackURL=/dashboard"
    })

When the OAuth provider returns the user to your server, the plugin automatically redirects them to the intended callback URL.

To share cookies between the proxy server and your main server it uses URL query parameters to pass the cookies encrypted in the URL. This is secure as the cookies are encrypted and can only be decrypted by the server.

This plugin requires skipping the state cookie check. This has security implications and should only be used in dev or staging environments. If `baseURL` and `productionURL` are the same, the plugin will not proxy the request.

**currentURL**: The application's current URL is automatically determined by the plugin. It first checks for the request URL if invoked by a client, then it checks the base URL from popular hosting providers, and finally falls back to the `baseURL` in your auth config. If the URL isn’t inferred correctly, you can specify it manually here.

**productionURL**: If this value matches the `baseURL` in your auth config, requests will not be proxied. Defaults to the `BETTER_AUTH_URL` environment variable.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/oauth-proxy.mdx)</content>
</page>

<page>
  <title>User & Accounts | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/users-accounts</url>
  <content>Beyond authenticating users, Better Auth also provides a set of methods to manage users. This includes, updating user information, changing passwords, and more.

The user table stores the authentication data of the user [Click here to view the schema](https://www.better-auth.com/docs/concepts/database#user).

The user table can be extended using [additional fields](https://www.better-auth.com/docs/concepts/database#extending-core-schema) or by plugins to store additional data.

### [Update User Information](#update-user-information)

To update user information, you can use the `updateUser` function provided by the client. The `updateUser` function takes an object with the following properties:

    await authClient.updateUser({
        image: "https://example.com/image.jpg",
        name: "John Doe",
    })

### [Change Email](#change-email)

To allow users to change their email, first enable the `changeEmail` feature, which is disabled by default. Set `changeEmail.enabled` to `true`:

    export const auth = betterAuth({
        user: {
            changeEmail: {
                enabled: true,
            }
        },
        emailVerification: {
            // Required to send the verification email
            sendVerificationEmail: async ({ user, url, token }) => {
                void sendEmail({
                    to: user.email,
                })
            }
        }
    })

Avoid awaiting the email sending to prevent timing attacks. On serverless platforms, use `waitUntil` or similar to ensure the email is sent.

By default, when a user requests to change their email, a verification email is sent to the **new** email address. The email is only updated after the user verifies the new email.

#### [Confirming with Current Email](#confirming-with-current-email)

For added security, you can require users to confirm the change via their **current** email before the verification email is sent to the new address. To do this, provide the `sendChangeEmailConfirmation` function.

    export const auth = betterAuth({
        user: {
            changeEmail: {
                enabled: true,
                sendChangeEmailConfirmation: async ({ user, newEmail, url, token }, request) => { 
                    void sendEmail({
                        to: user.email, // Sent to the CURRENT email
                        subject: 'Approve email change',
                        text: `Click the link to approve the change to ${newEmail}: ${url}`
                    })
                }
            }
        },
        // ...
    })

#### [Updating Without Verification](#updating-without-verification)

If you want to allow users to update their email immediately without verification (only if their current email is NOT verified), you can enable `updateEmailWithoutVerification`.

    export const auth = betterAuth({
        user: {
            changeEmail: {
                enabled: true,
                updateEmailWithoutVerification: true
            }
        }
    })

If `updateEmailWithoutVerification` is false (default), the email will not be updated until the new email is verified, even if the current email is unverified.

#### [Client Usage](#client-usage)

Use the `changeEmail` function on the client to initiate the process.

    await authClient.changeEmail({
        newEmail: "new-email@email.com",
        callbackURL: "/dashboard", // to redirect after verification
    });

### [Change Password](#change-password)

A user's password isn't stored in the user table. Instead, it's stored in the account table. To change the password of a user, you can use one of the following approaches:

    const { data, error } = await authClient.changePassword({    newPassword: "newpassword1234", // required    currentPassword: "oldpassword1234", // required    revokeOtherSessions: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `newPassword` | 
The new password to set

 | `string` |
| `currentPassword` | 

The current user password

 | `string` |
| `revokeOtherSessions?` | 

When set to true, all other active sessions for this user will be invalidated

 | `boolean` |

### [Set Password](#set-password)

If a user was registered using OAuth or other providers, they won't have a password or a credential account. In this case, you can use the `setPassword` action to set a password for the user. For security reasons, this function can only be called from the server. We recommend having users go through a 'forgot password' flow to set a password for their account.

    await auth.api.setPassword({
        body: { newPassword: "password" },
        headers: // headers containing the user's session token
    });

Better Auth provides a utility to hard delete a user from your database. It's disabled by default, but you can enable it easily by passing `enabled:true`

    export const auth = betterAuth({
        //...other config
        user: {
            deleteUser: { 
                enabled: true
            } 
        }
    })

Once enabled, you can call `authClient.deleteUser` to permanently delete user data from your database.

### [Adding Verification Before Deletion](#adding-verification-before-deletion)

For added security, you’ll likely want to confirm the user’s intent before deleting their account. A common approach is to send a verification email. Better Auth provides a `sendDeleteAccountVerification` utility for this purpose. This is especially needed if you have OAuth setup and want them to be able to delete their account without forcing them to login again for a fresh session.

Here’s how you can set it up:

    export const auth = betterAuth({
        user: {
            deleteUser: {
                enabled: true,
                sendDeleteAccountVerification: async (
                    {
                        user,   // The user object
                        url, // The auto-generated URL for deletion
                        token  // The verification token  (can be used to generate custom URL)
                    },
                    request  // The original request object (optional)
                ) => {
                    // Your email sending logic here
                    // Example: sendEmail(data.user.email, "Verify Deletion", data.url);
                },
            },
        },
    });

**How callback verification works:**

*   **Callback URL**: The URL provided in `sendDeleteAccountVerification` is a pre-generated link that deletes the user data when accessed.

delete-user.ts

    await authClient.deleteUser({
        callbackURL: "/goodbye" // you can provide a callback URL to redirect after deletion
    });

*   **Authentication Check**: The user must be signed in to the account they’re attempting to delete. If they aren’t signed in, the deletion process will fail.

If you have sent a custom URL, you can use the `deleteUser` method with the token to delete the user.

delete-user.ts

    await authClient.deleteUser({
        token
    });

### [Authentication Requirements](#authentication-requirements)

To delete a user, the user must meet one of the following requirements:

1.  A valid password

if the user has a password, they can delete their account by providing the password.

delete-user.ts

    await authClient.deleteUser({
        password: "password"
    });

2.  Fresh session

The user must have a `fresh` session token, meaning the user must have signed in recently. This is checked if the password is not provided.

By default `session.freshAge` is set to `60 * 60 * 24` (1 day). You can change this value by passing the `session` object to the `auth` configuration. If it is set to `0`, the freshness check is disabled. It is recommended not to disable this check if you are not using email verification for deleting the account.

delete-user.ts

    await authClient.deleteUser();

3.  Enabled email verification (needed for OAuth users)

As OAuth users don't have a password, we need to send a verification email to confirm the user's intent to delete their account. If you have already added the `sendDeleteAccountVerification` callback, you can just call the `deleteUser` method without providing any other information.

delete-user.ts

    await authClient.deleteUser();

4.  If you have a custom delete account page and sent that url via the `sendDeleteAccountVerification` callback. Then you need to call the `deleteUser` method with the token to complete the deletion.

delete-user.ts

    await authClient.deleteUser({
        token
    });

### [Callbacks](#callbacks)

**beforeDelete**: This callback is called before the user is deleted. You can use this callback to perform any cleanup or additional checks before deleting the user.

auth.ts

    export const auth = betterAuth({
        user: {
            deleteUser: {
                enabled: true,
                beforeDelete: async (user) => {
                    // Perform any cleanup or additional checks here
                },
            },
        },
    });

you can also throw `APIError` to interrupt the deletion process.

auth.ts

    import { betterAuth } from "better-auth";
    import { APIError } from "better-auth/api";
    
    export const auth = betterAuth({
        user: {
            deleteUser: {
                enabled: true,
                beforeDelete: async (user, request) => {
                    if (user.email.includes("admin")) {
                        throw new APIError("BAD_REQUEST", {
                            message: "Admin accounts can't be deleted",
                        });
                    }
                },
            },
        },
    });

**afterDelete**: This callback is called after the user is deleted. You can use this callback to perform any cleanup or additional actions after the user is deleted.

auth.ts

    export const auth = betterAuth({
        user: {
            deleteUser: {
                enabled: true,
                afterDelete: async (user, request) => {
                    // Perform any cleanup or additional actions here
                },
            },
        },
    });

Better Auth supports multiple authentication methods. Each authentication method is called a provider. For example, email and password authentication is a provider, Google authentication is a provider, etc.

When a user signs in using a provider, an account is created for the user. The account stores the authentication data returned by the provider. This data includes the access token, refresh token, and other information returned by the provider.

The account table stores the authentication data of the user [Click here to view the schema](https://www.better-auth.com/docs/concepts/database#account)

### [List User Accounts](#list-user-accounts)

To list user accounts you can use `client.user.listAccounts` method. Which will return all accounts associated with a user.

    const accounts = await authClient.listAccounts();

### [Token Encryption](#token-encryption)

Better Auth doesn’t encrypt tokens by default and that’s intentional. We want you to have full control over how encryption and decryption are handled, rather than baking in behavior that could be confusing or limiting. If you need to store encrypted tokens (like accessToken or refreshToken), you can use databaseHooks to encrypt them before they’re saved to your database.

    export const auth = betterAuth({
        databaseHooks: {
            account: {
                create: {
                    before(account, context) {
                        const withEncryptedTokens = { ...account };
                        if (account.accessToken) {
                            const encryptedAccessToken = encrypt(account.accessToken)  
                            withEncryptedTokens.accessToken = encryptedAccessToken;
                        }
                        if (account.refreshToken) {
                            const encryptedRefreshToken = encrypt(account.refreshToken); 
                            withEncryptedTokens.refreshToken = encryptedRefreshToken;
                        }
                        return {
                            data: withEncryptedTokens
                        }
                    },
                }
            }
        }
    })

Then whenever you retrieve back the account make sure to decrypt the tokens before using them.

### [Account Linking](#account-linking)

Account linking enables users to associate multiple authentication methods with a single account. With Better Auth, users can connect additional social sign-ons or OAuth providers to their existing accounts if the provider confirms the user's email as verified.

If account linking is disabled, no accounts can be linked, regardless of the provider or email verification status.

auth.ts

    export const auth = betterAuth({
        account: {
            accountLinking: {
                enabled: true, 
            }
        },
    });

#### [Forced Linking](#forced-linking)

You can specify a list of "trusted providers." When a user logs in using a trusted provider, their account will be automatically linked even if the provider doesn’t confirm the email verification status. Use this with caution as it may increase the risk of account takeover.

auth.ts

    export const auth = betterAuth({
        account: {
            accountLinking: {
                enabled: true,
                trustedProviders: ["google", "github"]
            }
        },
    });

#### [Manually Linking Accounts](#manually-linking-accounts)

Users already signed in can manually link their account to additional social providers or credential-based accounts.

*   **Linking Social Accounts:** Use the `linkSocial` method on the client to link a social provider to the user's account.
    
        await authClient.linkSocial({
            provider: "google", // Provider to link
            callbackURL: "/callback" // Callback URL after linking completes
        });
    
    You can also request specific scopes when linking a social account, which can be different from the scopes used during the initial authentication:
    
        await authClient.linkSocial({
            provider: "google",
            callbackURL: "/callback",
            scopes: ["https://www.googleapis.com/auth/drive.readonly"] // Request additional scopes
        });
    
    You can also link accounts using ID tokens directly, without redirecting to the provider's OAuth flow:
    
        await authClient.linkSocial({
            provider: "google",
            idToken: {
                token: "id_token_from_provider",
                nonce: "nonce_used_for_token", // Optional
                accessToken: "access_token", // Optional, may be required by some providers
                refreshToken: "refresh_token" // Optional
            }
        });
    
    This is useful when you already have valid tokens from the provider, for example:
    
    *   After signing in with a native SDK
    *   When using a mobile app that handles authentication
    *   When implementing custom OAuth flows
    
    The ID token must be valid and the provider must support ID token verification.
    
    If you want your users to be able to link a social account with a different email address than the user, or if you want to use a provider that does not return email addresses, you will need to enable this in the account linking settings.
    
    auth.ts
    
        export const auth = betterAuth({
            account: {
                accountLinking: {
                    allowDifferentEmails: true
                }
            },
        });
    
    If you want the newly linked accounts to update the user information, you need to enable this in the account linking settings.
    
    auth.ts
    
        export const auth = betterAuth({
            account: {
                accountLinking: {
                    updateUserInfoOnLink: true
                }
            },
        });
    
*   **Linking Credential-Based Accounts:** To link a credential-based account (e.g., email and password), users can initiate a "forgot password" flow, or you can call the `setPassword` method on the server.
    
        await auth.api.setPassword({
            headers: /* headers containing the user's session token */,
            password: /* new password */
        });
    

`setPassword` can't be called from the client for security reasons.

### [Account Unlinking](#account-unlinking)

You can unlink a user account by providing a `providerId`.

    await authClient.unlinkAccount({
        providerId: "google"
    });
    
    // Unlink a specific account
    await authClient.unlinkAccount({
        providerId: "google",
        accountId: "123"
    });

If the account doesn't exist, it will throw an error. Additionally, if the user only has one account, unlinking will be prevented to stop account lockout (unless `allowUnlinkingAll` is set to `true`).

auth.ts

    export const auth = betterAuth({
        account: {
            accountLinking: {
                allowUnlinkingAll: true
            }
        },
    });</content>
</page>

<page>
  <title>SvelteKit Example | Better Auth</title>
  <url>https://www.better-auth.com/docs/examples/svelte-kit</url>
  <content>This is an example of how to use Better Auth with SvelteKit.

**Implements the following features:** Email & Password . Social Sign-in with Google . Passkeys . Email Verification . Password Reset . Two Factor Authentication . Profile Update . Session Management

1.  Clone the code sandbox (or the repo) and open it in your code editor
2.  Move .env.example to .env and provide necessary variables
3.  Run the following commands
    
        pnpm install
        pnpm dev
    
4.  Open the browser and navigate to `http://localhost:3000`

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/examples/svelte-kit.mdx)</content>
</page>

<page>
  <title>Optimizing for Performance | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/optimizing-for-performance</url>
  <content>In this guide, we’ll go over some of the ways you can optimize your application for a more performant Better Auth app.

Caching is a powerful technique that can significantly improve the performance of your Better Auth application by reducing the number of database queries and speeding up response times.

### [Cookie Cache](#cookie-cache)

Calling your database every time `useSession` or `getSession` is invoked isn’t ideal, especially if sessions don’t change frequently. Cookie caching handles this by storing session data in a short-lived, signed cookie similar to how JWT access tokens are used with refresh tokens.

To turn on cookie caching, just set `session.cookieCache` in your auth config:

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      session: {
        cookieCache: {
          enabled: true,
          maxAge: 5 * 60, // Cache duration in seconds
        },
      },
    });

Read more about [cookie caching](https://www.better-auth.com/docs/concepts/session-management#cookie-cache).

### [Framework Caching](#framework-caching)

Here are examples of how you can do caching in different frameworks and environments:

If you're using a framework that supports server-side rendering, it's usually best to pre-fetch the user session on the server and use it as a fallback on the client.

    const session = await auth.api.getSession({
      headers: await headers(),
    });
    //then pass the session to the client

Optimizing database performance is essential to get the best out of Better Auth.

#### [Recommended fields to index](#recommended-fields-to-index)

| Table | Fields | Plugin |
| --- | --- | --- |
| users | `email` |  |
| accounts | `userId` |  |
| sessions | `userId`, `token` |  |
| verifications | `identifier` |  |
| invitations | `email`, `organizationId` | organization |
| members | `userId`, `organizationId` | organization |
| organizations | `slug` | organization |
| passkey | `userId` | passkey |
| twoFactor | `secret` | twoFactor |

We intend to add indexing support in our schema generation tool in the future.

If you're using custom adapters (like Prisma, Drizzle, or MongoDB), you can reduce your bundle size by using `better-auth/minimal` instead of `better-auth`. This version excludes Kysely, which is only needed when using direct database connections.

### [Usage](#usage)

Simply import from `better-auth/minimal` instead of `better-auth`:

**Limitations:**

*   Direct database connections are not supported (you must use an adapter)
*   Built-in migrations are not supported. Use external migration tools (or use `better-auth` if you need built-in migration support)

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/guides/optimizing-for-performance.mdx)</content>
</page>

<page>
  <title>Astro Example | Better Auth</title>
  <url>https://www.better-auth.com/docs/examples/astro</url>
  <content>This is an example of how to use Better Auth with Astro. It uses Solid for building the components.

**Implements the following features:** Email & Password . Social Sign-in with Google . Passkeys . Email Verification . Password Reset . Two Factor Authentication . Profile Update . Session Management

1.  Clone the code sandbox (or the repo) and open it in your code editor
    
2.  Provide .env file with the following variables
    
        GOOGLE_CLIENT_ID=
        GOOGLE_CLIENT_SECRET=
        BETTER_AUTH_SECRET=
    
    //if you don't have these, you can get them from the google developer console. If you don't want to use google sign-in, you can remove the google config from the `auth.ts` file.
    
3.  Run the following commands
    
        pnpm install
        pnpm run dev
    
4.  Open the browser and navigate to `http://localhost:3000`
    

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/examples/astro.mdx)</content>
</page>

<page>
  <title>Other Social Providers | Better Auth</title>
  <url>https://www.better-auth.com/docs/authentication/other-social-providers</url>
  <content>Better Auth provides support for any social provider that implements the OAuth2 protocol or OpenID Connect (OIDC) flows through the [Generic OAuth Plugin](https://www.better-auth.com/docs/plugins/generic-oauth). You can use pre-configured helper functions for popular providers like Auth0, Keycloak, Okta, Microsoft Entra ID, and Slack, or manually configure any OAuth provider.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

To use the Generic OAuth plugin, add it to your auth config.

auth.ts

    import { betterAuth } from "better-auth"
    import { genericOAuth } from "better-auth/plugins"
    
    export const auth = betterAuth({
        // ... other config options
        plugins: [
            genericOAuth({ 
                config: [ 
                    { 
                        providerId: "provider-id", 
                        clientId: "test-client-id", 
                        clientSecret: "test-client-secret", 
                        discoveryUrl: "https://auth.example.com/.well-known/openid-configuration", 
                        // ... other config options
                    }, 
                    // Add more providers as needed
                ] 
            }) 
        ]
    })

### [Add the client plugin](#add-the-client-plugin)

Include the Generic OAuth client plugin in your authentication client instance.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { genericOAuthClient } from "better-auth/client/plugins"
    
    const authClient = createAuthClient({
        plugins: [
            genericOAuthClient()
        ]
    })

Read more about installation and usage of the Generic Oauth plugin [here](https://www.better-auth.com/docs/plugins/generic-oauth#usage).

Here's a basic example of configuring a generic OAuth provider:

auth.ts

    import { betterAuth } from "better-auth"
    import { genericOAuth } from "better-auth/plugins"
    
    export const auth = betterAuth({
      plugins: [
        genericOAuth({
          config: [
            {
              providerId: "provider-id",
              clientId: process.env.CLIENT_ID,
              clientSecret: process.env.CLIENT_SECRET,
              discoveryUrl: "https://auth.example.com/.well-known/openid-configuration",
            },
          ],
        }),
      ],
    })

Better Auth provides pre-configured helper functions for popular OAuth providers. Here's an example using Slack:

auth.ts

    import { betterAuth } from "better-auth"
    import { genericOAuth, slack } from "better-auth/plugins"
    
    export const auth = betterAuth({
      plugins: [
        genericOAuth({
          config: [
            slack({
              clientId: process.env.SLACK_CLIENT_ID,
              clientSecret: process.env.SLACK_CLIENT_SECRET,
            }),
          ],
        }),
      ],
    })

sign-in.ts

    const response = await authClient.signIn.oauth2({
      providerId: "slack",
      callbackURL: "/dashboard",
    })

For more pre-configured providers (Auth0, Keycloak, Okta, Microsoft Entra ID) and their configuration options, see the [Generic OAuth Plugin documentation](https://www.better-auth.com/docs/plugins/generic-oauth#pre-configured-provider-helpers).

If you need to configure a provider that doesn't have a pre-configured helper, you can manually configure it:

### [Instagram Example](#instagram-example)

auth.ts

    import { betterAuth } from "better-auth";
    import { genericOAuth } from "better-auth/plugins";
    
    export const auth = betterAuth({
      // ... other config options
      plugins: [
        genericOAuth({
          config: [
            {
              providerId: "instagram",
              clientId: process.env.INSTAGRAM_CLIENT_ID as string,
              clientSecret: process.env.INSTAGRAM_CLIENT_SECRET as string,
              authorizationUrl: "https://api.instagram.com/oauth/authorize",
              tokenUrl: "https://api.instagram.com/oauth/access_token",
              scopes: ["user_profile", "user_media"],
            },
          ],
        }),
      ],
    });

sign-in.ts

    const response = await authClient.signIn.oauth2({
      providerId: "instagram",
      callbackURL: "/dashboard", // the path to redirect to after the user is authenticated
    });

### [Coinbase Example](#coinbase-example)

auth.ts

    import { betterAuth } from "better-auth";
    import { genericOAuth } from "better-auth/plugins";
    
    export const auth = betterAuth({
      // ... other config options
      plugins: [
        genericOAuth({
          config: [
            {
              providerId: "coinbase",
              clientId: process.env.COINBASE_CLIENT_ID as string,
              clientSecret: process.env.COINBASE_CLIENT_SECRET as string,
              authorizationUrl: "https://www.coinbase.com/oauth/authorize",
              tokenUrl: "https://api.coinbase.com/oauth/token",
              scopes: ["wallet:user:read"], // and more...
            },
          ],
        }),
      ],
    });

sign-in.ts

    const response = await authClient.signIn.oauth2({
      providerId: "coinbase",
      callbackURL: "/dashboard", // the path to redirect to after the user is authenticated
    });

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/authentication/other-social-providers.mdx)</content>
</page>

<page>
  <title>OAuth | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/oauth</url>
  <content>Better Auth comes with built-in support for OAuth 2.0 and OpenID Connect. This allows you to authenticate users via popular OAuth providers like Google, Facebook, GitHub, and more.

If your desired provider isn't directly supported, you can use the [Generic OAuth Plugin](https://www.better-auth.com/docs/plugins/generic-oauth) for custom integrations.

To enable a social provider, you need to provide `clientId` and `clientSecret` for the provider.

Here's an example of how to configure Google as a provider:

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
        },
      },
    });

### [Sign In](#sign-in)

To sign in with a social provider, you can use the `signIn.social` function with the `authClient` or `auth.api` for server-side usage.

    // client-side usage
    await authClient.signIn.social({
      provider: "google", // or any other provider id
    })

    // server-side usage
    await auth.api.signInSocial({
      body: {
        provider: "google", // or any other provider id
      },
    });

### [Link account](#link-account)

To link an account to a social provider, you can use the `linkAccount` function with the `authClient` or `auth.api` for server-side usage.

    await authClient.linkSocial({
      provider: "google", // or any other provider id
    })

server-side usage:

    await auth.api.linkSocialAccount({
      body: {
        provider: "google", // or any other provider id
      },
      headers: // pass headers with authenticated token
    });

### [Get Access Token](#get-access-token)

To get the access token for a social provider, you can use the `getAccessToken` function with the `authClient` or `auth.api` for server-side usage. When you use this endpoint, if the access token is expired, it will be refreshed.

    const { accessToken } = await authClient.getAccessToken({
      providerId: "google", // or any other provider id
      accountId: "accountId", // optional, if you want to get the access token for a specific account
    })

server-side usage:

    await auth.api.getAccessToken({
      body: {
        providerId: "google", // or any other provider id
        accountId: "accountId", // optional, if you want to get the access token for a specific account
        userId: "userId", // optional, if you don't provide headers with authenticated token
      },
      headers: // pass headers with authenticated token
    });

### [Get Account Info Provided by the provider](#get-account-info-provided-by-the-provider)

To get provider specific account info you can use the `accountInfo` function with the `authClient` or `auth.api` for server-side usage.

    const info = await authClient.accountInfo({
      accountId: "accountId", // here you pass in the provider given account id, the provider is automatically detected from the account id
    })

server-side usage:

    await auth.api.accountInfo({
      body: { accountId: "accountId" },
      headers: // pass headers with authenticated token
    });

### [Requesting Additional Scopes](#requesting-additional-scopes)

Sometimes your application may need additional OAuth scopes after the user has already signed up (e.g., for accessing GitHub repositories or Google Drive). Users may not want to grant extensive permissions initially, preferring to start with minimal permissions and grant additional access as needed.

You can request additional scopes by using the `linkSocial` method with the same provider. This will trigger a new OAuth flow that requests the additional scopes while maintaining the existing account connection.

    const requestAdditionalScopes = async () => {
        await authClient.linkSocial({
            provider: "google",
            scopes: ["https://www.googleapis.com/auth/drive.file"],
        });
    };

Make sure you're running Better Auth version 1.2.7 or later. Earlier versions (like 1.2.2) may show a "Social account already linked" error when trying to link with an existing provider for additional scopes.

### [Passing Additional Data Through OAuth Flow](#passing-additional-data-through-oauth-flow)

Better Auth allows you to pass additional data through the OAuth flow without storing it in the database. This is useful for scenarios like tracking referral codes, analytics sources, or other temporary data that should be processed during authentication but not persisted.

When initiating OAuth sign-in or account linking, pass the additional data:

    // Client-side: Sign in with additional data
    await authClient.signIn.social({
      provider: "google",
      additionalData: {
        referralCode: "ABC123",
        source: "landing-page",
      },
    });
    
    // Client-side: Link account with additional data
    await authClient.linkSocial({
      provider: "google",
      additionalData: {
        referralCode: "ABC123",
      },
    });
    
    // Server-side: Sign in with additional data
    await auth.api.signInSocial({
      body: {
        provider: "google",
        additionalData: {
          referralCode: "ABC123",
          source: "admin-panel",
        },
      },
    });

#### [Accessing Additional Data in Hooks](#accessing-additional-data-in-hooks)

The additional data is available in your hooks during the OAuth callback through the `getOAuthState`.

This usually works for `/callback/:id` paths and the generic OAuth plugin callback path (`/oauth2/callback/:providerId`).

Example using a before hook:

auth.ts

    import { betterAuth } from "better-auth";
    import { getOAuthState } from "better-auth/api";
    
    export const auth = betterAuth({
      // Other configurations...
      hooks: {
        after: [
          {
            matcher: () => true,
            handler: async (ctx) => {
              // Additional data is only available during OAuth callback
              if (ctx.path === "/callback/:id") {
                const additionalData = await getOAuthState<{
                  referralCode?: string;
                  source?: string;
                }>();
    
                if (additionalData) {
                  // IMPORTANT: Validate and sanitize the data before using it
                  // This data comes from the client and should not be trusted
    
                  // Example: Validate and process referral code
                  if (additionalData.referralCode) {
                    const isValidFormat = /^[A-Z0-9]{6}$/.test(additionalData.referralCode);
                    if (isValidFormat) {
                      // Verify the referral code exists in your database
                      const referral = await db.referrals.findByCode(additionalData.referralCode);
                      if (referral) {
                        // Safe to use the verified referral
                        await db.referrals.incrementUsage(referral.id);
                      }
                    }
                  }
    
                  // Track analytics (low-risk usage)
                  if (additionalData.source) {
                    await analytics.track("oauth_signin", {
                      source: additionalData.source,
                      userId: ctx.context.session?.user.id,
                    });
                  }
                }
              }
            },
          },
        ],
      },
    });

Example using a database hook:

auth.ts

     // You can also access additional data in database hooks
      databaseHooks: {
        user: {
          create: {
            before: async (user, ctx) => {
              if (ctx.path === "/callback/:id") {
                const additionalData = await getOAuthState<{ referredFrom?: string }>();
                if (additionalData?.referredFrom) {
                  return {
                    data: {
                      referredFrom: additionalData.referredFrom,
                    },
                  };
                }
              }
            },
          },
        },
      },

By default OAuth state includes the following data:

*   `callbackURL` - the callback URL for the OAuth flow
*   `codeVerifier` - the code verifier for the OAuth flow
*   `errorURL` - the error URL for the OAuth flow
*   `newUserURL` - the new user URL for the OAuth flow
*   `link` - the link for the OAuth flow (email and user id)
*   `requestSignUp` - whether to request sign up for the OAuth flow
*   `expiresAt` - the expiration time of the OAuth state
*   `[key: string]`: any additional data you pass in the OAuth flow

### [scope](#scope)

The scope of the access request. For example, `email` or `profile`.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          scope: ["email", "profile"],
        },
      },
    });

### [redirectURI](#redirecturi)

Custom redirect URI for the provider. By default, it uses `/api/auth/callback/${providerName}`

auth.ts

    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          redirectURI: "https://your-app.com/auth/callback",
        },
      },
    });

### [disableSignUp](#disablesignup)

Disables sign-up for new users.

### [disableIdTokenSignIn](#disableidtokensignin)

Disables the use of the ID token for sign-in. By default, it's enabled for some providers like Google and Apple.

### [verifyIdToken](#verifyidtoken)

A custom function to verify the ID token.

### [overrideUserInfoOnSignIn](#overrideuserinfoonsignin)

A boolean value that determines whether to override the user information in the database when signing in. By default, it is set to `false`, meaning that the user information will not be overridden during sign-in. If you want to update the user information every time they sign in, set this to `true`.

### [mapProfileToUser](#mapprofiletouser)

A custom function to map the user profile returned from the provider to the user object in your database.

Useful, if you have additional fields in your user object you want to populate from the provider's profile. Or if you want to change how by default the user object is mapped.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          mapProfileToUser: (profile) => {
            return {
              firstName: profile.given_name,
              lastName: profile.family_name,
            };
          },
        },
      },
    });

### [refreshAccessToken](#refreshaccesstoken)

A custom function to refresh the token. This feature is only supported for built-in social providers (Google, Facebook, GitHub, etc.) and is not currently supported for custom OAuth providers configured through the Generic OAuth Plugin. For built-in providers, you can provide a custom function to refresh the token if needed.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          refreshAccessToken: async (token) => {
            return {
              accessToken: "new-access-token",
              refreshToken: "new-refresh-token",
            };
          },
        },
      },
    });

### [clientKey](#clientkey)

The client key of your application. This is used by TikTok Social Provider instead of `clientId`.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        tiktok: {
          clientKey: "YOUR_TIKTOK_CLIENT_KEY",
          clientSecret: "YOUR_TIKTOK_CLIENT_SECRET",
        },
      },
    });

### [getUserInfo](#getuserinfo)

A custom function to get user info from the provider. This allows you to override the default user info retrieval process.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          getUserInfo: async (token) => {
            // Custom implementation to get user info
            const response = await fetch("https://www.googleapis.com/oauth2/v2/userinfo", {
              headers: {
                Authorization: `Bearer ${token.accessToken}`,
              },
            });
            const profile = await response.json();
            return {
              user: {
                id: profile.id,
                name: profile.name,
                email: profile.email,
                image: profile.picture,
                emailVerified: profile.verified_email,
              },
              data: profile,
            };
          },
        },
      },
    });

### [disableImplicitSignUp](#disableimplicitsignup)

Disables implicit sign up for new users. When set to true for the provider, sign-in needs to be called with `requestSignUp` as true to create new users.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          disableImplicitSignUp: true,
        },
      },
    });

### [prompt](#prompt)

The prompt to use for the authorization code request. This controls the authentication flow behavior.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          prompt: "select_account", // or "consent", "login", "none", "select_account+consent"
        },
      },
    });

### [responseMode](#responsemode)

The response mode to use for the authorization code request. This determines how the authorization response is returned.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          responseMode: "query", // or "form_post"
        },
      },
    });

### [disableDefaultScope](#disabledefaultscope)

Removes the default scopes of the provider. By default, providers include certain scopes like `email` and `profile`. Set this to `true` to remove these default scopes and use only the scopes you specify.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      // Other configurations...
      socialProviders: {
        google: {
          clientId: "YOUR_GOOGLE_CLIENT_ID",
          clientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
          disableDefaultScope: true,
          scope: ["https://www.googleapis.com/auth/userinfo.email"], // Only this scope will be used
        },
      },
    });

### [Other Provider Configurations](#other-provider-configurations)

Each provider may have additional options, check the specific provider documentation for more details.</content>
</page>

<page>
  <title>Remix Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/remix</url>
  <content>Better Auth can be easily integrated with Remix. This guide will show you how to integrate Better Auth with Remix.

You can follow the steps from [installation](https://www.better-auth.com/docs/installation) to get started or you can follow this guide to make it the Remix-way.

If you have followed the installation steps, you can skip the first step.

[Create auth instance](#create-auth-instance)
---------------------------------------------

Create a file named `auth.server.ts` in one of these locations:

*   Project root
*   `lib/` folder
*   `utils/` folder

You can also nest any of these folders under `app/` folder. (e.g. `app/lib/auth.server.ts`)

And in this file, import Better Auth and create your instance.

Make sure to export the auth instance with the variable name `auth` or as a `default` export.

app/lib/auth.server.ts

    import { betterAuth } from "better-auth"
    
    export const auth = betterAuth({
        database: {
            provider: "postgres", //change this to your database provider
            url: process.env.DATABASE_URL, // path to your database or connection string
        }
    })

[Create API Route](#create-api-route)
-------------------------------------

We need to mount the handler to a API route. Create a resource route file `api.auth.$.ts` inside `app/routes/` directory. And add the following code:

app/routes/api.auth.$.ts

    import { auth } from '~/lib/auth.server' // Adjust the path as necessary
    import type { LoaderFunctionArgs, ActionFunctionArgs } from "@remix-run/node"
    
    export async function loader({ request }: LoaderFunctionArgs) {
        return auth.handler(request)
    }
    
    export async function action({ request }: ActionFunctionArgs) {
        return auth.handler(request)
    }

You can change the path on your better-auth configuration but it's recommended to keep it as `routes/api.auth.$.ts`

[Create a client](#create-a-client)
-----------------------------------

Create a client instance. Here we are creating `auth-client.ts` file inside the `lib/` directory.

app/lib/auth-client.ts

    import { createAuthClient } from "better-auth/react" // make sure to import from better-auth/react
    
    export const authClient = createAuthClient({
        //you can pass client configuration here
    })

Once you have created the client, you can use it to sign up, sign in, and perform other actions.

### [Example usage](#example-usage)

#### [Sign Up](#sign-up)

app/routes/signup.tsx

    import { Form } from "@remix-run/react"
    import { useState } from "react"
    import { authClient } from "~/lib/auth-client"
    
    export default function SignUp() {
      const [email, setEmail] = useState("")
      const [name, setName] = useState("")
      const [password, setPassword] = useState("")
    
      const signUp = async () => {
        await authClient.signUp.email(
          {
            email,
            password,
            name,
          },
          {
            onRequest: (ctx) => {
              // show loading state
            },
            onSuccess: (ctx) => {
              // redirect to home
            },
            onError: (ctx) => {
              alert(ctx.error)
            },
          },
        )
      }
    
      return (
        <div>
          <h2>
            Sign Up
          </h2>
          <Form
            onSubmit={signUp}
          >
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Name"
            />
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Email"
            />
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Password"
            />
            <button
              type="submit"
            >
              Sign Up
            </button>
          </Form>
        </div>
      )
    }

#### [Sign In](#sign-in)

app/routes/signin.tsx

    import { Form } from "@remix-run/react"
    import { useState } from "react"
    import { authClient } from "~/services/auth-client"
    
    export default function SignIn() {
      const [email, setEmail] = useState("")
      const [password, setPassword] = useState("")
    
      const signIn = async () => {
        await authClient.signIn.email(
          {
            email,
            password,
          },
          {
            onRequest: (ctx) => {
              // show loading state
            },
            onSuccess: (ctx) => {
              // redirect to home
            },
            onError: (ctx) => {
              alert(ctx.error)
            },
          },
        )
      }
    
      return (
        <div>
          <h2>
            Sign In
          </h2>
          <Form onSubmit={signIn}>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <button
              type="submit"
            >
              Sign In
            </button>
          </Form>
        </div>
      )
    }</content>
</page>

<page>
  <title>Plugins | Better Auth</title>
  <url>https://www.better-auth.com/docs/concepts/plugins</url>
  <content>Plugins are a key part of Better Auth, they let you extend the base functionalities. You can use them to add new authentication methods, features, or customize behaviors.

Better Auth comes with many built-in plugins ready to use. Check the plugins section for details. You can also create your own plugins.

Plugins can be a server-side plugin, a client-side plugin, or both.

To add a plugin on the server, include it in the `plugins` array in your auth configuration. The plugin will initialize with the provided options.

server.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        plugins: [
            // Add your plugins here
        ]
    });

Client plugins are added when creating the client. Most plugins require both server and client plugins to work correctly. The Better Auth auth client on the frontend uses the `createAuthClient` function provided by `better-auth/client`.

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    
    const authClient = createAuthClient({
        plugins: [
            // Add your client plugins here
        ]
    });

We recommend keeping the auth-client and your normal auth instance in separate files.

To get started, you'll need a server plugin. Server plugins are the backbone of all plugins, and client plugins are there to provide an interface with frontend APIs to easily work with your server plugins.

If your server plugins have endpoints that need to be called from the client, you'll also need to create a client plugin.

### [What can a plugin do?](#what-can-a-plugin-do)

*   Create custom `endpoint`s to perform any action you want.
*   Extend database tables with custom `schemas`.
*   Use a `middleware` to target a group of routes using its route matcher, and run only when those routes are called through a request.
*   Use `hooks` to target a specific route or request. And if you want to run the hook even if the endpoint is called directly.
*   Use `onRequest` or `onResponse` if you want to do something that affects all requests or responses.
*   Create a custom `rate-limit` rule.

To create a server plugin, you need to pass an object that satisfies the `BetterAuthPlugin` interface.

The only required property is `id`, which is a unique identifier for the plugin. Both server and client plugins can use the same `id`.

plugin.ts

    import type { BetterAuthPlugin } from "better-auth";
    
    export const myPlugin = () => {
        return {
            id: "my-plugin",
        } satisfies BetterAuthPlugin
    }

You don't have to make the plugin a function, but it's recommended to do so. This way, you can pass options to the plugin and it's consistent with the built-in plugins.

### [Endpoints](#endpoints)

To add endpoints to the server, you can pass `endpoints` which requires an object with the key being any `string` and the value being an `AuthEndpoint`.

To create an Auth Endpoint you'll need to import `createAuthEndpoint` from `better-auth`.

Better Auth uses wraps around another library called [Better Call](https://github.com/bekacru/better-call) to create endpoints. Better call is a simple ts web framework made by the same team behind Better Auth.

plugin.ts

    import { createAuthEndpoint } from "better-auth/api";
    
    const myPlugin = () => {
        return {
            id: "my-plugin",
            endpoints: {
                getHelloWorld: createAuthEndpoint("/my-plugin/hello-world", {
                    method: "GET",
                }, async(ctx) => {
                    return ctx.json({
                        message: "Hello World"
                    })
                })
            }
        } satisfies BetterAuthPlugin
    }

Create Auth endpoints wraps around `createEndpoint` from Better Call. Inside the `ctx` object, it'll provide another object called `context` that give you access better-auth specific contexts including `options`, `db`, `baseURL` and more.

**Context Object**

*   `appName`: The name of the application. Defaults to "Better Auth".
*   `options`: The options passed to the Better Auth instance.
*   `tables`: Core tables definition. It is an object which has the table name as the key and the schema definition as the value.
*   `baseURL`: the baseURL of the auth server. This includes the path. For example, if the server is running on `http://localhost:3000`, the baseURL will be `http://localhost:3000/api/auth` by default unless changed by the user.
*   `session`: The session configuration. Includes `updateAge` and `expiresIn` values.
*   `secret`: The secret key used for various purposes. This is defined by the user.
*   `authCookie`: The default cookie configuration for core auth cookies.
*   `logger`: The logger instance used by Better Auth.
*   `db`: The Kysely instance used by Better Auth to interact with the database.
*   `adapter`: This is the same as db but it give you `orm` like functions to interact with the database. (we recommend using this over `db` unless you need raw sql queries or for performance reasons)
*   `internalAdapter`: These are internal db calls that are used by Better Auth. For example, you can use these calls to create a session instead of using `adapter` directly. `internalAdapter.createSession(userId)`
*   `createAuthCookie`: This is a helper function that lets you get a cookie `name` and `options` for either to `set` or `get` cookies. It implements things like `__Secure-` prefix for cookies based on whether the connection is secure (HTTPS) or the application is running in production mode.
*   `trustedOrigins`: This is the list of trusted origins that you specified via `options.trustedOrigins`.
*   `isTrustedOrigin`: This is a helper function that allows you to quickly check whether a given url or path is trusted based on the trusted origins configuration.

For other properties, you can check the [Better Call](https://github.com/bekacru/better-call) documentation and the [source code](https://github.com/better-auth/better-auth/blob/main/packages/better-auth/src/init.ts) .

**Rules for Endpoints**

*   Makes sure you use kebab-case for the endpoint path
*   Make sure to only use `POST` or `GET` methods for the endpoints.
*   Any function that modifies a data should be a `POST` method.
*   Any function that fetches data should be a `GET` method.
*   Make sure to use the `createAuthEndpoint` function to create API endpoints.
*   Make sure your paths are unique to avoid conflicts with other plugins. If you're using a common path, add the plugin name as a prefix to the path. (`/my-plugin/hello-world` instead of `/hello-world`.)

### [Schema](#schema)

You can define a database schema for your plugin by passing a `schema` object. The schema object should have the table name as the key and the schema definition as the value.

plugin.ts

    import { BetterAuthPlugin } from "better-auth/plugins";
    
    const myPlugin = () => {
        return {
            id: "my-plugin",
            schema: {
                myTable: {
                    fields: {
                        name: {
                            type: "string"
                        }
                    },
                    modelName: "myTable" // optional if you want to use a different name than the key
                }
            }
        } satisfies BetterAuthPlugin
    }

**Fields**

By default better-auth will create an `id` field for each table. You can add additional fields to the table by adding them to the `fields` object.

The key is the column name and the value is the column definition. The column definition can have the following properties:

`type`: The type of the field. It can be `string`, `number`, `boolean`, `date`.

`required`: if the field should be required on a new record. (default: `false`)

`unique`: if the field should be unique. (default: `false`)

`reference`: if the field is a reference to another table. (default: `null`) It takes an object with the following properties:

*   `model`: The table name to reference.
*   `field`: The field name to reference.
*   `onDelete`: The action to take when the referenced record is deleted. (default: `null`)

**Other Schema Properties**

`disableMigration`: if the table should not be migrated. (default: `false`)

plugin.ts

    const myPlugin = (opts: PluginOptions) => {
        return {
            id: "my-plugin",
            schema: {
                rateLimit: {
                    fields: {
                        key: {
                            type: "string",
                        },
                    },
                    disableMigration: opts.storage.provider !== "database", 
                },
            },
        } satisfies BetterAuthPlugin
    }

if you add additional fields to a `user` or `session` table, the types will be inferred automatically on `getSession` and `signUpEmail` calls.

plugin.ts

    
    const myPlugin = () => {
        return {
            id: "my-plugin",
            schema: {
                user: {
                    fields: {
                        age: {
                            type: "number",
                        },
                    },
                },
            },
        } satisfies BetterAuthPlugin
    }

This will add an `age` field to the `user` table and all `user` returning endpoints will include the `age` field and it'll be inferred properly by typescript.

Don't store sensitive information in the `user` or `session` table. Create a new table if you need to store sensitive information.

### [Hooks](#hooks)

Hooks are used to run code before or after an action is performed, either from a client or directly on the server. You can add hooks to the server by passing a `hooks` object, which should contain `before` and `after` properties.

plugin.ts

    import { createAuthMiddleware } from "better-auth/plugins";
    
    const myPlugin = () => {
        return {
            id: "my-plugin",
            hooks: {
                before: [{
                        matcher: (context) => {
                            return context.headers.get("x-my-header") === "my-value"
                        },
                        handler: createAuthMiddleware(async (ctx) => {
                            // do something before the request
                            return  {
                                context: ctx // if you want to modify the context
                            }
                        })
                    }],
                after: [{
                    matcher: (context) => {
                        return context.path === "/sign-up/email"
                    },
                    handler: createAuthMiddleware(async (ctx) => {
                        return ctx.json({
                            message: "Hello World"
                        }) // if you want to modify the response
                    })
                }]
            }
        } satisfies BetterAuthPlugin
    }

### [Middleware](#middleware)

You can add middleware to the server by passing a `middlewares` array. This array should contain middleware objects, each with a `path` and a `middleware` property. Unlike hooks, middleware only runs on `api` requests from a client. If the endpoint is invoked directly, the middleware will not run.

The `path` can be either a string or a path matcher, using the same path-matching system as `better-call`.

If you throw an `APIError` from the middleware or return a `Response` object, the request will be stopped, and the response will be sent to the client.

plugin.ts

    const myPlugin = () => {
        return {
            id: "my-plugin",
            middlewares: [
                {
                    path: "/my-plugin/hello-world",
                    middleware: createAuthMiddleware(async(ctx) => {
                        // do something
                    })
                }
            ]
        } satisfies BetterAuthPlugin
    }

### [On Request & On Response](#on-request--on-response)

Additional to middlewares, you can also hook into right before a request is made and right after a response is returned. This is mostly useful if you want to do something that affects all requests or responses.

#### [On Request](#on-request)

The `onRequest` function is called right before the request is made. It takes two parameters: the `request` and the `context` object.

Here’s how it works:

*   **Continue as Normal**: If you don't return anything, the request will proceed as usual.
*   **Interrupt the Request**: To stop the request and send a response, return an object with a `response` property that contains a `Response` object.
*   **Modify the Request**: You can also return a modified `request` object to change the request before it's sent.

plugin.ts

    const myPlugin = () => {
        return  {
            id: "my-plugin",
            onRequest: async (request, context) => {
                // do something
            },
        } satisfies BetterAuthPlugin
    }

#### [On Response](#on-response)

The `onResponse` function is executed immediately after a response is returned. It takes two parameters: the `response` and the `context` object.

Here’s how to use it:

*   **Modify the Response**: You can return a modified response object to change the response before it is sent to the client.
*   **Continue Normally**: If you don't return anything, the response will be sent as is.

plugin.ts

    const myPlugin = () => {
        return {
            id: "my-plugin",
            onResponse: async (response, context) => {
                // do something
            },
        } satisfies BetterAuthPlugin
    }

### [Rate Limit](#rate-limit)

You can define custom rate limit rules for your plugin by passing a `rateLimit` array. The rate limit array should contain an array of rate limit objects.

plugin.ts

    const myPlugin = () => {
        return {
            id: "my-plugin",
            rateLimit: [
                {
                    pathMatcher: (path) => {
                        return path === "/my-plugin/hello-world"
                    },
                    limit: 10,
                    window: 60,
                }
            ]
        } satisfies BetterAuthPlugin
    }

### [Trusted origins](#trusted-origins)

If you're building custom plugins or endpoints, you can use the `isTrustedOrigin()` method available on the auth context to validate URLs against your trusted origins configuration. This ensures your custom endpoints respect the same security settings as Better Auth's built-in endpoints.

plugin.ts

    import { createAuthEndpoint, APIError } from "better-auth/api";
    import * as z from "zod"
    
    const myPlugin = () => {
        return {
            id: "my-plugin",
            trustedOrigins: [
                "http://trusted.com"
            ],
            endpoints: {
                getTrustedHelloWorld: createAuthEndpoint("/my-plugin/hello-world", {
                    method: "GET",
                    query: z.object({
                        url: z.string()
                    }),
                }, async (ctx) => {
                    // The allowRelativePaths option can be used to either allow or disallow relative paths
                    if (!ctx.context.isTrustedOrigin(ctx.query.url, { allowRelativePaths: false })) {
                        throw new APIError("FORBIDDEN", {
                            message: "origin is not trusted."
                        });
                    }
    
                    return ctx.json({
                        message: "Hello World"
                    })
                })
            }
        } satisfies BetterAuthPlugin
    }

See the [trusted origins and security](https://www.better-auth.com/docs/reference/security#trusted-origins) docs for more info.

### [Server-plugin helper functions](#server-plugin-helper-functions)

Some additional helper functions for creating server plugins.

#### [`getSessionFromCtx`](#getsessionfromctx)

Allows you to get the client's session data by passing the auth middleware's `context`.

plugin.ts

    import { createAuthMiddleware } from "better-auth/plugins";
    import { getSessionFromCtx } from "better-auth/api";
    
    const myPlugin = {
        id: "my-plugin",
        hooks: {
            before: [{
                    matcher: (context) => {
                        return context.headers.get("x-my-header") === "my-value"
                    },
                    handler: createAuthMiddleware(async (ctx) => {
                        const session = await getSessionFromCtx(ctx);
                        // do something with the client's session.
    
                        return  {
                            context: ctx
                        }
                    })
                }],
        }
    } satisfies BetterAuthPlugin

#### [`sessionMiddleware`](#sessionmiddleware)

A middleware that checks if the client has a valid session. If the client has a valid session, it'll add the session data to the context object.

plugin.ts

    import { createAuthMiddleware } from "better-auth/plugins";
    import { sessionMiddleware } from "better-auth/api";
    
    const myPlugin = () => {
        return {
            id: "my-plugin",
            endpoints: {
                getHelloWorld: createAuthEndpoint("/my-plugin/hello-world", {
                    method: "GET",
                    use: [sessionMiddleware], 
                }, async (ctx) => {
                    const session = ctx.context.session;
                    return ctx.json({
                        message: "Hello World"
                    })
                })
            }
        } satisfies BetterAuthPlugin
    }

If your endpoints need to be called from the client, you'll also need to create a client plugin. Better Auth clients can infer the endpoints from the server plugins. You can also add additional client-side logic.

client-plugin.ts

    import type { BetterAuthClientPlugin } from "better-auth";
    
    export const myPluginClient = () => {
        return {
            id: "my-plugin",
        } satisfies BetterAuthClientPlugin
    }

### [Endpoint Interface](#endpoint-interface)

Endpoints are inferred from the server plugin by adding a `$InferServerPlugin` key to the client plugin.

The client infers the `path` as an object and converts kebab-case to camelCase. For example, `/my-plugin/hello-world` becomes `myPlugin.helloWorld`.

client-plugin.ts

    import type { BetterAuthClientPlugin } from "better-auth/client";
    import type { myPlugin } from "./plugin";
    
    const myPluginClient = () => {
        return  {
            id: "my-plugin",
            $InferServerPlugin: {} as ReturnType<typeof myPlugin>,
        } satisfies BetterAuthClientPlugin
    }

### [Get actions](#get-actions)

If you need to add additional methods or whatnot to the client, you can use the `getActions` function. This function is called with the `fetch` function from the client.

Better Auth uses [Better fetch](https://better-fetch.vercel.app/) to make requests. Better Fetch is a simple fetch wrapper made by the same author of Better Auth.

client-plugin.ts

    import type { BetterAuthClientPlugin } from "better-auth/client";
    import type { myPlugin } from "./plugin";
    import type { BetterFetchOption } from "@better-fetch/fetch";
    
    const myPluginClient = {
        id: "my-plugin",
        $InferServerPlugin: {} as ReturnType<typeof myPlugin>,
        getActions: ($fetch) => {
            return {
                myCustomAction: async (data: {
                    foo: string,
                }, fetchOptions?: BetterFetchOption) => {
                    const res = $fetch("/custom/action", {
                        method: "POST",
                        body: {
                            foo: data.foo
                        },
                        ...fetchOptions
                    })
                    return res
                }
            }
        }
    } satisfies BetterAuthClientPlugin

As a general guideline, ensure that each function accepts only one argument, with an optional second argument for fetchOptions to allow users to pass additional options to the fetch call. The function should return an object containing data and error keys.

If your use case involves actions beyond API calls, feel free to deviate from this rule.

### [Get Atoms](#get-atoms)

This is only useful if you want to provide `hooks` like `useSession`.

Get atoms is called with the `fetch` function from better fetch, and it should return an object with the atoms. The atoms should be created using [nanostores](https://github.com/nanostores/nanostores). The atoms will be resolved by each framework's `useStore` hook provided by nanostores.

client-plugin.ts

    import { atom } from "nanostores";
    import type { BetterAuthClientPlugin } from "better-auth/client";
    
    const myPluginClient = {
        id: "my-plugin",
        $InferServerPlugin: {} as ReturnType<typeof myPlugin>,
        getAtoms: ($fetch) => {
            const myAtom = atom<null>()
            return {
                myAtom
            }
        }
    } satisfies BetterAuthClientPlugin

See built-in plugins for examples of how to use atoms properly.

### [Path methods](#path-methods)

By default, inferred paths use the `GET` method if they don't require a body and `POST` if they do. You can override this by passing a `pathMethods` object. The key should be the path, and the value should be the method ("POST" | "GET").

client-plugin.ts

    import type { BetterAuthClientPlugin } from "better-auth/client";
    import type { myPlugin } from "./plugin";
    
    const myPluginClient = {
        id: "my-plugin",
        $InferServerPlugin: {} as ReturnType<typeof myPlugin>,
        pathMethods: {
            "/my-plugin/hello-world": "POST"
        }
    } satisfies BetterAuthClientPlugin

### [Fetch plugins](#fetch-plugins)

If you need to use better fetch plugins, you can pass them to the `fetchPlugins` array. You can read more about better fetch plugins in the [better fetch documentation](https://better-fetch.vercel.app/docs/plugins).

### [Atom Listeners](#atom-listeners)

This is only useful if you want to provide `hooks` like `useSession` and you want to listen to atoms and re-evaluate them when they change.

You can see how this is used in the built-in plugins.</content>
</page>

<page>
  <title>MongoDB Adapter | Better Auth</title>
  <url>https://www.better-auth.com/docs/adapters/mongo</url>
  <content>MongoDB is a popular NoSQL database that is widely used for building scalable and flexible applications. It provides a flexible schema that allows for easy data modeling and querying.

Before getting started, make sure you have MongoDB installed and configured. For more information, see [MongoDB Documentation](https://www.mongodb.com/docs/)

You can use the MongoDB adapter to connect to your database as follows.

auth.ts

    import { betterAuth } from "better-auth";
    import { MongoClient } from "mongodb";
    import { mongodbAdapter } from "better-auth/adapters/mongodb";
    
    const client = new MongoClient("mongodb://localhost:27017/database");
    const db = client.db();
    
    export const auth = betterAuth({
      database: mongodbAdapter(db, {
        // Optional: if you don't provide a client, database transactions won't be enabled.
        client
      }),
    });

For MongoDB, we don't need to generate or migrate the schema.

Database joins is useful when Better-Auth needs to fetch related data from multiple tables in a single query. Endpoints like `/get-session`, `/get-full-organization` and many others benefit greatly from this feature, seeing upwards of 2x to 3x performance improvements depending on database latency.

The MongoDB adapter supports joins out of the box since version `1.4.0`. To enable this feature, you need to set the `experimental.joins` option to `true` in your auth configuration.

auth.ts

    export const auth = betterAuth({
      experimental: { joins: true }
    });

*   If you're looking for performance improvements or tips, take a look at our guide to [performance optimizations](https://www.better-auth.com/docs/guides/optimizing-for-performance).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/adapters/mongo.mdx)</content>
</page>

<page>
  <title>Email & Password | Better Auth</title>
  <url>https://www.better-auth.com/docs/authentication/email-password#forget-password</url>
  <content>Email and password authentication is a common method used by many applications. Better Auth provides a built-in email and password authenticator that you can easily integrate into your project.

If you prefer username-based authentication, check out the [username plugin](https://www.better-auth.com/docs/plugins/username). It extends the email and password authenticator with username support.

To enable email and password authentication, you need to set the `emailAndPassword.enabled` option to `true` in the `auth` configuration.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
      emailAndPassword: { 
        enabled: true, 
      }, 
    });

If it's not enabled, it'll not allow you to sign in or sign up with email and password.

### [Sign Up](#sign-up)

To sign a user up, you can use the `signUp.email` function provided by the client.

    const { data, error } = await authClient.signUp.email({    name: "John Doe", // required    email: "john.doe@example.com", // required    password: "password1234", // required    image: "https://example.com/image.png",    callbackURL: "https://example.com/callback",});

| Prop | Description | Type |
| --- | --- | --- |
| `name` | 
The name of the user.

 | `string` |
| `email` | 

The email address of the user.

 | `string` |
| `password` | 

The password of the user. It should be at least 8 characters long and max 128 by default.

 | `string` |
| `image?` | 

An optional profile image of the user.

 | `string` |
| `callbackURL?` | 

An optional URL to redirect to after the user signs up.

 | `string` |

These are the default properties for the sign up email endpoint, however it's possible that with [additional fields](https://www.better-auth.com/docs/concepts/typescript#additional-fields) or special plugins you can pass more properties to the endpoint.

### [Sign In](#sign-in)

To sign a user in, you can use the `signIn.email` function provided by the client.

    const { data, error } = await authClient.signIn.email({    email: "john.doe@example.com", // required    password: "password1234", // required    rememberMe: true,    callbackURL: "https://example.com/callback",});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
The email address of the user.

 | `string` |
| `password` | 

The password of the user. It should be at least 8 characters long and max 128 by default.

 | `string` |
| `rememberMe?` | 

If false, the user will be signed out when the browser is closed. (optional) (default: true)

 | `boolean` |
| `callbackURL?` | 

An optional URL to redirect to after the user signs in. (optional)

 | `string` |

These are the default properties for the sign in email endpoint, however it's possible that with [additional fields](https://www.better-auth.com/docs/concepts/typescript#additional-fields) or special plugins you can pass different properties to the endpoint.

### [Sign Out](#sign-out)

To sign a user out, you can use the `signOut` function provided by the client.

    await authClient.signOut();

you can pass `fetchOptions` to redirect onSuccess

auth-client.ts

    await authClient.signOut({
      fetchOptions: {
        onSuccess: () => {
          router.push("/login"); // redirect to login page
        },
      },
    });

### [Email Verification](#email-verification)

To enable email verification, you need to pass a function that sends a verification email with a link. The `sendVerificationEmail` function takes a data object with the following properties:

*   `user`: The user object.
*   `url`: The URL to send to the user which contains the token.
*   `token`: A verification token used to complete the email verification.

and a `request` object as the second parameter.

auth.ts

    import { betterAuth } from "better-auth";
    import { sendEmail } from "./email"; // your email sending function
    
    export const auth = betterAuth({
      emailVerification: {
        sendVerificationEmail: async ( { user, url, token }, request) => {
          void sendEmail({
            to: user.email,
            subject: "Verify your email address",
            text: `Click the link to verify your email: ${url}`,
          });
        },
      },
    });

Avoid awaiting the email sending to prevent timing attacks. On serverless platforms, use `waitUntil` or similar to ensure the email is sent.

On the client side you can use `sendVerificationEmail` function to send verification link to user. This will trigger the `sendVerificationEmail` function you provided in the `auth` configuration.

Once the user clicks on the link in the email, if the token is valid, the user will be redirected to the URL provided in the `callbackURL` parameter. If the token is invalid, the user will be redirected to the URL provided in the `callbackURL` parameter with an error message in the query string `?error=invalid_token`.

#### [Require Email Verification](#require-email-verification)

If you enable require email verification, users must verify their email before they can log in. And every time a user tries to sign in, sendVerificationEmail is called.

This only works if you have sendVerificationEmail implemented and if the user is trying to sign in with email and password.

auth.ts

    export const auth = betterAuth({
      emailAndPassword: {
        requireEmailVerification: true,
      },
    });

If a user tries to sign in without verifying their email, you can handle the error and show a message to the user.

auth-client.ts

    await authClient.signIn.email(
      {
        email: "email@example.com",
        password: "password",
      },
      {
        onError: (ctx) => {
          // Handle the error
          if (ctx.error.status === 403) {
            alert("Please verify your email address");
          }
          //you can also show the original error message
          alert(ctx.error.message);
        },
      }
    );

#### [Triggering manually Email Verification](#triggering-manually-email-verification)

You can trigger the email verification manually by calling the `sendVerificationEmail` function.

    await authClient.sendVerificationEmail({
      email: "user@email.com",
      callbackURL: "/", // The redirect URL after verification
    });

### [Request Password Reset](#request-password-reset)

To allow users to reset a password first you need to provide `sendResetPassword` function to the email and password authenticator. The `sendResetPassword` function takes a data object with the following properties:

*   `user`: The user object.
*   `url`: The URL to send to the user which contains the token.
*   `token`: A verification token used to complete the password reset.

and a `request` object as the second parameter.

auth.ts

    import { betterAuth } from "better-auth";
    import { sendEmail } from "./email"; // your email sending function
    
    export const auth = betterAuth({
      emailAndPassword: {
        enabled: true,
        sendResetPassword: async ({user, url, token}, request) => {
          void sendEmail({
            to: user.email,
            subject: "Reset your password",
            text: `Click the link to reset your password: ${url}`,
          });
        },
        onPasswordReset: async ({ user }, request) => {
          // your logic here
          console.log(`Password for user ${user.email} has been reset.`);
        },
      },
    });

Avoid awaiting the email sending to prevent timing attacks. On serverless platforms, use `waitUntil` or similar to ensure the email is sent.

Additionally, you can provide an `onPasswordReset` callback to execute logic after a password has been successfully reset.

Once you configured your server you can call `requestPasswordReset` function to send reset password link to user. If the user exists, it will trigger the `sendResetPassword` function you provided in the auth config.

POST

/request-password-reset

    const { data, error } = await authClient.requestPasswordReset({    email: "john.doe@example.com", // required    redirectTo: "https://example.com/reset-password",});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
The email address of the user to send a password reset email to

 | `string` |
| `redirectTo?` | 

The URL to redirect the user to reset their password. If the token isn't valid or expired, it'll be redirected with a query parameter `?error=INVALID_TOKEN`. If the token is valid, it'll be redirected with a query parameter \`?token=VALID\_TOKEN

 | `string` |

When a user clicks on the link in the email, they will be redirected to the reset password page. You can add the reset password page to your app. Then you can use `resetPassword` function to reset the password. It takes an object with the following properties:

*   `newPassword`: The new password of the user.

auth-client.ts

    const { data, error } = await authClient.resetPassword({
      newPassword: "password1234",
      token,
    });

    const token = new URLSearchParams(window.location.search).get("token");if (!token) {  // Handle the error}const { data, error } = await authClient.resetPassword({    newPassword: "password1234", // required    token, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `newPassword` | 
The new password to set

 | `string` |
| `token` | 

The token to reset the password

 | `string` |

### [Update password](#update-password)

A user's password isn't stored in the user table. Instead, it's stored in the account table. To change the password of a user, you can use one of the following approaches:

    const { data, error } = await authClient.changePassword({    newPassword: "newpassword1234", // required    currentPassword: "oldpassword1234", // required    revokeOtherSessions: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `newPassword` | 
The new password to set

 | `string` |
| `currentPassword` | 

The current user password

 | `string` |
| `revokeOtherSessions?` | 

When set to true, all other active sessions for this user will be invalidated

 | `boolean` |

### [Configuration](#configuration)

**Password**

Better Auth stores passwords inside the `account` table with `providerId` set to `credential`.

**Password Hashing**: Better Auth uses `scrypt` to hash passwords. The `scrypt` algorithm is designed to be slow and memory-intensive to make it difficult for attackers to brute force passwords. OWASP recommends using `scrypt` if `argon2id` is not available. We decided to use `scrypt` because it's natively supported by Node.js.

You can pass custom password hashing algorithm by setting `passwordHasher` option in the `auth` configuration.

auth.ts

    import { betterAuth } from "better-auth"
    import { scrypt } from "scrypt"
    
    export const auth = betterAuth({
        //...rest of the options
        emailAndPassword: {
            password: {
                hash: // your custom password hashing function
                verify: // your custom password verification function
            }
        }
    })</content>
</page>

<page>
  <title>SvelteKit Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/svelte-kit</url>
  <content>Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Mount the handler](#mount-the-handler)

We need to mount the handler to SvelteKit server hook.

hooks.server.ts

    import { auth } from "$lib/auth";
    import { svelteKitHandler } from "better-auth/svelte-kit";
    import { building } from "$app/environment";
    
    export async function handle({ event, resolve }) {
      return svelteKitHandler({ event, resolve, auth, building });
    }

### [Populate session data in the event (`event.locals`)](#populate-session-data-in-the-event-eventlocals)

The `svelteKitHandler` does not automatically populate `event.locals.user` or `event.locals.session`. If you want to access the current session in your server code (e.g., in `+layout.server.ts`, actions, or endpoints), populate `event.locals` in your `handle` hook:

hooks.server.ts

    import { auth } from "$lib/auth";
    import { svelteKitHandler } from "better-auth/svelte-kit";
    import { building } from "$app/environment";
    
    export async function handle({ event, resolve }) {
      // Fetch current session from Better Auth
      const session = await auth.api.getSession({
        headers: event.request.headers,
      });
    
      // Make session and user available on server
      if (session) {
        event.locals.session = session.session;
        event.locals.user = session.user;
      }
    
      return svelteKitHandler({ event, resolve, auth, building });
    }

### [Server Action Cookies](#server-action-cookies)

To ensure cookies are properly set when you call functions like `signInEmail` or `signUpEmail` in a server action, you should use the `sveltekitCookies` plugin. This plugin will automatically handle setting cookies for you in SvelteKit.

You need to add it as a plugin to your Better Auth instance.

The `getRequestEvent` function is available in SvelteKit `2.20.0` and later. Make sure you are using a compatible version.

lib/auth.ts

    import { betterAuth } from "better-auth";
    import { sveltekitCookies } from "better-auth/svelte-kit";
    import { getRequestEvent } from "$app/server";
    
    export const auth = betterAuth({
      // ... your config
      plugins: [sveltekitCookies(getRequestEvent)], // make sure this is the last plugin in the array
    });

Create a client instance. You can name the file anything you want. Here we are creating `client.ts` file inside the `lib/` directory.

auth-client.ts

    import { createAuthClient } from "better-auth/svelte"; // make sure to import from better-auth/svelte
    
    export const authClient = createAuthClient({
      // you can pass client configuration here
    });

Once you have created the client, you can use it to sign up, sign in, and perform other actions. Some of the actions are reactive. The client use [nano-store](https://github.com/nanostores/nanostores) to store the state and reflect changes when there is a change like a user signing in or out affecting the session state.

### [Example usage](#example-usage)

    <script lang="ts">
      import { authClient } from "$lib/client";
      const session = authClient.useSession();
    </script>
        <div>
          {#if $session.data}
            <div>
              <p>
                {$session.data.user.name}
              </p>
              <button
                on:click={async () => {
                  await authClient.signOut();
                }}
              >
                Sign Out
              </button>
            </div>
          {:else}
            <button
              on:click={async () => {
                await authClient.signIn.social({
                  provider: "github",
                });
              }}
            >
              Continue with GitHub
            </button>
          {/if}
        </div>

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/svelte-kit.mdx)</content>
</page>

<page>
  <title>Elysia Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/elysia</url>
  <content>This integration guide is assuming you are using Elysia with bun server.

Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Mount the handler](#mount-the-handler)

We need to mount the handler to Elysia endpoint.

    import { Elysia } from "elysia";
    import { auth } from "./auth";
    
    const app = new Elysia().mount(auth.handler).listen(3000);
    
    console.log(
      `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`,
    );

### [CORS](#cors)

To configure cors, you can use the `cors` plugin from `@elysiajs/cors`.

    import { Elysia } from "elysia";
    import { cors } from "@elysiajs/cors";
    
    import { auth } from "./auth";
    
    const app = new Elysia()
      .use(
        cors({
          origin: "http://localhost:3001",
          methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
          credentials: true,
          allowedHeaders: ["Content-Type", "Authorization"],
        }),
      )
      .mount(auth.handler)
      .listen(3000);
    
    console.log(
      `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`,
    );

### [Macro](#macro)

You can use [macro](https://elysiajs.com/patterns/macro.html#macro) with [resolve](https://elysiajs.com/essential/handler.html#resolve) to provide session and user information before pass to view.

    import { Elysia } from "elysia";
    import { auth } from "./auth";
    
    // user middleware (compute user and session and pass to routes)
    const betterAuth = new Elysia({ name: "better-auth" })
      .mount(auth.handler)
      .macro({
        auth: {
          async resolve({ status, request: { headers } }) {
            const session = await auth.api.getSession({
              headers,
            });
    
            if (!session) return status(401);
    
            return {
              user: session.user,
              session: session.session,
            };
          },
        },
      });
    
    const app = new Elysia()
      .use(betterAuth)
      .get("/user", ({ user }) => user, {
        auth: true,
      })
      .listen(3000);
    
    console.log(
      `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`,
    );

This will allow you to access the `user` and `session` object in all of your routes.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/elysia.mdx)</content>
</page>

<page>
  <title>One Tap | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/one-tap</url>
  <content>The One Tap plugin allows users to log in with a single tap using Google's One Tap API. The plugin provides a simple way to integrate One Tap into your application, handling the client-side and server-side logic for you.

### [Add the Server Plugin](#add-the-server-plugin)

Add the One Tap plugin to your auth configuration:

auth.ts

    import { betterAuth } from "better-auth";
    import { oneTap } from "better-auth/plugins"; 
    
    export const auth = betterAuth({
        plugins: [ 
            oneTap(), // Add the One Tap server plugin
        ] 
    });

### [Add the Client Plugin](#add-the-client-plugin)

Add the client plugin and specify where the user should be redirected after sign-in or if additional verification (like 2FA) is needed.

    import { createAuthClient } from "better-auth/client";
    import { oneTapClient } from "better-auth/client/plugins";
    
    export const authClient = createAuthClient({
      plugins: [
        oneTapClient({
          clientId: "YOUR_CLIENT_ID",
          // Optional client configuration:
          autoSelect: false,
          cancelOnTapOutside: true,
          context: "signin",
          additionalOptions: {
            // Any extra options for the Google initialize method
          },
          // Configure prompt behavior and exponential backoff:
          promptOptions: {
            baseDelay: 1000,   // Base delay in ms (default: 1000)
            maxAttempts: 5     // Maximum number of attempts before triggering onPromptNotification (default: 5)
          }
        })
      ]
    });

### [Usage](#usage)

To display the One Tap popup, simply call the oneTap method on your auth client:

    await authClient.oneTap();

### [Customizing Redirect Behavior](#customizing-redirect-behavior)

By default, after a successful login the plugin will hard redirect the user to `/`. You can customize this behavior as follows:

#### [Avoiding a Hard Redirect](#avoiding-a-hard-redirect)

Pass fetchOptions with an onSuccess callback to handle the login response without a page reload:

    await authClient.oneTap({
      fetchOptions: {
        onSuccess: () => {
          // For example, use a router to navigate without a full reload:
          router.push("/dashboard");
        }
      }
    });

#### [Specifying a Custom Callback URL](#specifying-a-custom-callback-url)

To perform a hard redirect to a different page after login, use the callbackURL option:

    await authClient.oneTap({
      callbackURL: "/dashboard"
    });

#### [Handling Prompt Dismissals with Exponential Backoff](#handling-prompt-dismissals-with-exponential-backoff)

If the user dismisses or skips the prompt, the plugin will retry showing the One Tap prompt using exponential backoff based on your configured promptOptions.

If the maximum number of attempts is reached without a successful sign-in, you can use the onPromptNotification callback to be notified—allowing you to render an alternative UI (e.g., a traditional Google Sign-In button) so users can restart the process manually:

    await authClient.oneTap({
      onPromptNotification: (notification) => {
        console.warn("Prompt was dismissed or skipped. Consider displaying an alternative sign-in option.", notification);
        // Render your alternative UI here
      }
    });

### [Client Options](#client-options)

*   **clientId**: The client ID for your Google One Tap API.
*   **autoSelect**: Automatically select the account if the user is already signed in. Default is false.
*   **context**: The context in which the One Tap API should be used (e.g., "signin"). Default is "signin".
*   **cancelOnTapOutside**: Cancel the One Tap popup when the user taps outside it. Default is true.
*   additionalOptions: Extra options to pass to Google's initialize method as per the [Google Identity Services docs](https://developers.google.com/identity/gsi/web/reference/js-reference#google.accounts.id.prompt).
*   **promptOptions**: Configuration for the prompt behavior and exponential backoff:
*   **baseDelay**: Base delay in milliseconds for retries. Default is 1000.
*   **maxAttempts**: Maximum number of prompt attempts before invoking the onPromptNotification callback. Default is 5.
*   **fedCM**: Whether to enable [Federated Credential Management](https://developer.mozilla.org/en-US/docs/Web/API/FedCM_API) (FedCM) support. Default is true.

### [Server Options](#server-options)

*   **disableSignUp**: Disable the sign-up option, allowing only existing users to sign in. Default is `false`.
*   **ClientId**: Optionally, pass a client ID here if it is not provided in your social provider configuration.

Ensure you have configured the Authorized JavaScript origins (e.g., [http://localhost:3000](http://localhost:3000/), [https://example.com](https://example.com/)) for your Client ID in the Google Cloud Console. This is a required step for the Google One Tap API, and it will not function correctly unless your origins are correctly set.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/one-tap.mdx)</content>
</page>

<page>
  <title>One-Time Token Plugin | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/one-time-token</url>
  <content>The One-Time Token (OTT) plugin provides functionality to generate and verify secure, single-use session tokens. These are commonly used for across domains authentication.

[Installation](#installation)
-----------------------------

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

To use the One-Time Token plugin, add it to your auth config.

auth.ts

    import { betterAuth } from "better-auth";
    import { oneTimeToken } from "better-auth/plugins/one-time-token";
    
    export const auth = betterAuth({
        plugins: [
          oneTimeToken()
        ]
        // ... other auth config
    });

### [Add the client plugin](#add-the-client-plugin)

Next, include the one-time-token client plugin in your authentication client instance.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { oneTimeTokenClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            oneTimeTokenClient()
        ]
    })

[Usage](#usage)
---------------

### [1\. Generate a Token](#1-generate-a-token)

Generate a token using `auth.api.generateOneTimeToken` or `authClient.oneTimeToken.generate`

GET

/one-time-token/generate

    const { data, error } = await authClient.oneTimeToken.generate();

GET

/one-time-token/generate

    const data = await auth.api.generateOneTimeToken({    // This endpoint requires session cookies.    headers: await headers(),});

This will return a `token` that is attached to the current session which can be used to verify the one-time token. By default, the token will expire in 3 minutes.

### [2\. Verify the Token](#2-verify-the-token)

When the user clicks the link or submits the token, use the `auth.api.verifyOneTimeToken` or `authClient.oneTimeToken.verify` method in another API route to validate it.

POST

/one-time-token/verify

    const { data, error } = await authClient.oneTimeToken.verify({    token: "some-token", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `token` | 
The token to verify.

 | `string` |

POST

/one-time-token/verify

    const data = await auth.api.verifyOneTimeToken({    body: {        token: "some-token", // required    },});

| Prop | Description | Type |
| --- | --- | --- |
| `token` | 
The token to verify.

 | `string` |

This will return the session that was attached to the token.

[Options](#options)
-------------------

These options can be configured when adding the `oneTimeToken` plugin:

*   **`disableClientRequest`** (boolean): Optional. If `true`, the token will only be generated on the server side. Default: `false`.
*   **`expiresIn`** (number): Optional. The duration for which the token is valid in minutes. Default: `3`.

    oneTimeToken({
        expiresIn: 10 // 10 minutes
    })

*   **`generateToken`**: A custom token generator function that takes `session` object and a `ctx` as parameters.
    
*   **`storeToken`**: Optional. This option allows you to configure how the token is stored in your database.
    
    *   **`plain`**: The token is stored in plain text. (Default)
    *   **`hashed`**: The token is hashed using the default hasher.
    *   **`custom-hasher`**: A custom hasher function that takes a token and returns a hashed token.

Note: It will not affect the token that's sent, it will only affect the token stored in your database.

Examples:

No hashing (default)

    oneTimeToken({
        storeToken: "plain"
    })

built-in hasher

    oneTimeToken({
        storeToken: "hashed"
    })

custom hasher

    oneTimeToken({
        storeToken: {
            type: "custom-hasher",
            hash: async (token) => {
                return myCustomHasher(token);
            }
        }
    })</content>
</page>

<page>
  <title>Lynx Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/lynx</url>
  <content>This integration guide is for using Better Auth with [Lynx](https://lynxjs.org/), a cross-platform rendering framework that enables developers to build applications for Android, iOS, and Web platforms with native rendering performance.

Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

Install Better Auth and the Lynx React dependency:

Import `createAuthClient` from `better-auth/lynx` to create your client instance:

lib/auth-client.ts

    import { createAuthClient } from "better-auth/lynx"
    
    export const authClient = createAuthClient({
        baseURL: "http://localhost:3000" // The base URL of your auth server
    })

The Lynx client provides the same API as other Better Auth clients, with optimized integration for Lynx's reactive system.

### [Authentication Methods](#authentication-methods)

    import { authClient } from "./lib/auth-client"
    
    // Sign in with email and password
    await authClient.signIn.email({
        email: "test@user.com",
        password: "password1234"
    })
    
    // Sign up
    await authClient.signUp.email({
        email: "test@user.com", 
        password: "password1234",
        name: "John Doe"
    })
    
    // Sign out
    await authClient.signOut()

### [Hooks](#hooks)

The Lynx client includes reactive hooks that integrate seamlessly with Lynx's component system:

#### [useSession](#usesession)

components/user.tsx

    import { authClient } from "../lib/auth-client"
    
    export function User() {
        const {
            data: session,
            isPending, // loading state
            error // error object 
        } = authClient.useSession()
    
        if (isPending) return <div>Loading...</div>
        if (error) return <div>Error: {error.message}</div>
    
        return (
            <div>
                {session ? (
                    <div>
                        <p>Welcome, {session.user.name}!</p>
                        <button onClick={() => authClient.signOut()}>
                            Sign Out
                        </button>
                    </div>
                ) : (
                    <button onClick={() => authClient.signIn.social({
                        provider: 'github'
                    })}>
                        Sign In with GitHub
                    </button>
                )}
            </div>
        )
    }

### [Store Integration](#store-integration)

The Lynx client uses [nanostores](https://github.com/nanostores/nanostores) for state management and provides a `useStore` hook for accessing reactive state:

components/session-info.tsx

    import { useStore } from "better-auth/lynx"
    import { authClient } from "../lib/auth-client"
    
    export function SessionInfo() {
        // Access the session store directly
        const session = useStore(authClient.$store.session)
        
        return (
            <div>
                {session && (
                    <pre>{JSON.stringify(session, null, 2)}</pre>
                )}
            </div>
        )
    }

### [Advanced Store Usage](#advanced-store-usage)

You can use the store with selective key watching for optimized re-renders:

components/optimized-user.tsx

    import { useStore } from "better-auth/lynx"
    import { authClient } from "../lib/auth-client"
    
    export function OptimizedUser() {
        // Only re-render when specific keys change
        const session = useStore(authClient.$store.session, {
            keys: ['user.name', 'user.email'] // Only watch these specific keys
        })
        
        return (
            <div>
                {session?.user && (
                    <div>
                        <h2>{session.user.name}</h2>
                        <p>{session.user.email}</p>
                    </div>
                )}
            </div>
        )
    }

The Lynx client supports all Better Auth plugins:

lib/auth-client.ts

    import { createAuthClient } from "better-auth/lynx"
    import { magicLinkClient } from "better-auth/client/plugins"
    
    const authClient = createAuthClient({
        plugins: [
            magicLinkClient()
        ]
    })
    
    // Use plugin methods
    await authClient.signIn.magicLink({
        email: "test@email.com"
    })

Error handling works the same as other Better Auth clients:

components/login-form.tsx

    import { authClient } from "../lib/auth-client"
    
    export function LoginForm() {
        const signIn = async (email: string, password: string) => {
            const { data, error } = await authClient.signIn.email({
                email,
                password
            })
            
            if (error) {
                console.error('Login failed:', error.message)
                return
            }
            
            console.log('Login successful:', data)
        }
        
        return (
            <form onSubmit={(e) => {
                e.preventDefault()
                const formData = new FormData(e.target)
                signIn(formData.get('email'), formData.get('password'))
            }}>
                <input name="email" type="email" placeholder="Email" />
                <input name="password" type="password" placeholder="Password" />
                <button type="submit">Sign In</button>
            </form>
        )
    }

The Lynx client provides:

*   **Cross-Platform Support**: Works across Android, iOS, and Web platforms
*   **Optimized Performance**: Built specifically for Lynx's reactive system
*   **Nanostores Integration**: Uses nanostores for efficient state management
*   **Selective Re-rendering**: Watch specific store keys to minimize unnecessary updates
*   **Full API Compatibility**: All Better Auth methods and plugins work seamlessly
*   **TypeScript Support**: Full type safety with TypeScript inference

The Lynx integration maintains all the features and benefits of Better Auth while providing optimal performance and developer experience within Lynx's cross-platform ecosystem.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/lynx.mdx)</content>
</page>

<page>
  <title>Options | Better Auth</title>
  <url>https://www.better-auth.com/docs/reference/options</url>
  <content>List of all the available options for configuring Better Auth. See [Better Auth Options](https://github.com/better-auth/better-auth/blob/main/packages/core/src/types/init-options.ts).

The name of the application.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	appName: "My App",
    })

Base URL for Better Auth. This is typically the root URL where your application server is hosted. Note: If you include a path in the baseURL, it will take precedence over the default path.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	baseURL: "https://example.com",
    })

If not explicitly set, the system will check for the environment variable `process.env.BETTER_AUTH_URL`

Base path for Better Auth. This is typically the path where the Better Auth routes are mounted. It will be overridden if there is a path component within `baseURL`.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	basePath: "/api/auth",
    })

Default: `/api/auth`

List of trusted origins. You can provide a static array of origins, a function that returns origins dynamically, or use wildcard patterns to match multiple domains.

### [Static Origins](#static-origins)

You can provide a static array of origins:

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	trustedOrigins: ["http://localhost:3000", "https://example.com"],
    })

### [Dynamic Origins](#dynamic-origins)

You can provide a function that returns origins dynamically:

    export const auth = betterAuth({
    	trustedOrigins: async (request: Request) => {
    		// Return an array of trusted origins based on the request
    		return ["https://dynamic-origin.com"];
    	}
    })

### [Wildcard Support](#wildcard-support)

You can use wildcard patterns in trusted origins:

    export const auth = betterAuth({
    	trustedOrigins: [
    		"https://*.example.com", // trust all HTTPS subdomains of example.com
    		"http://*.dev.example.com" // trust all HTTP subdomains of dev.example.com
    	]
    })

Make sure to provide the protocol prefix when using wildcard patterns. For example, `https://*.example.com` instead of `*.example.com`.

The secret used for encryption, signing, and hashing.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	secret: "your-secret-key",
    })

By default, Better Auth will look for the following environment variables:

*   `process.env.BETTER_AUTH_SECRET`
*   `process.env.AUTH_SECRET`

If none of these environment variables are set, it will default to `"better-auth-secret-123456789"`. In production, if it's not set, it will throw an error.

You can generate a good secret using the following command:

    openssl rand -base64 32

Database configuration for Better Auth.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	database: {
    		dialect: "postgres",
    		type: "postgres",
    		casing: "camel"
    	},
    })

Better Auth supports various database configurations including [PostgreSQL](https://www.better-auth.com/docs/adapters/postgresql), [MySQL](https://www.better-auth.com/docs/adapters/mysql), and [SQLite](https://www.better-auth.com/docs/adapters/sqlite).

Read more about databases [here](https://www.better-auth.com/docs/concepts/database).

Secondary storage configuration used to store session and rate limit data.

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
    	// ... other options
        secondaryStorage: {
        	// Your implementation here
        },
    })

Read more about secondary storage [here](https://www.better-auth.com/docs/concepts/database#secondary-storage).

Email verification configuration.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	emailVerification: {
    		sendVerificationEmail: async ({ user, url, token }) => {
    			// Send verification email to user
    		},
    		sendOnSignUp: true,
    		autoSignInAfterVerification: true,
    		expiresIn: 3600 // 1 hour
    	},
    })

*   `sendVerificationEmail`: Function to send verification email
*   `sendOnSignUp`: Send verification email automatically after sign up (default: `false`)
*   `sendOnSignIn`: Send verification email automatically on sign in when the user's email is not verified (default: `false`)
*   `autoSignInAfterVerification`: Auto sign in the user after they verify their email
*   `expiresIn`: Number of seconds the verification token is valid for (default: `3600` seconds)

Email and password authentication configuration.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	emailAndPassword: {
    		enabled: true,
    		disableSignUp: false,
    		requireEmailVerification: true,
    		minPasswordLength: 8,
    		maxPasswordLength: 128,
    		autoSignIn: true,
    		sendResetPassword: async ({ user, url, token }) => {
    			// Send reset password email
    		},
    		resetPasswordTokenExpiresIn: 3600, // 1 hour
    		password: {
    			hash: async (password) => {
    				// Custom password hashing
    				return hashedPassword;
    			},
    			verify: async ({ hash, password }) => {
    				// Custom password verification
    				return isValid;
    			}
    		}
    	},
    })

*   `enabled`: Enable email and password authentication (default: `false`)
*   `disableSignUp`: Disable email and password sign up (default: `false`)
*   `requireEmailVerification`: Require email verification before a session can be created
*   `minPasswordLength`: Minimum password length (default: `8`)
*   `maxPasswordLength`: Maximum password length (default: `128`)
*   `autoSignIn`: Automatically sign in the user after sign up
*   `sendResetPassword`: Function to send reset password email
*   `resetPasswordTokenExpiresIn`: Number of seconds the reset password token is valid for (default: `3600` seconds)
*   `password`: Custom password hashing and verification functions

Configure social login providers.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	socialProviders: {
    		google: {
    			clientId: "your-client-id",
    			clientSecret: "your-client-secret",
    			redirectURI: "https://example.com/api/auth/callback/google"
    		},
    		github: {
    			clientId: "your-client-id",
    			clientSecret: "your-client-secret",
    			redirectURI: "https://example.com/api/auth/callback/github"
    		}
    	},
    })

List of Better Auth plugins.

    import { betterAuth } from "better-auth";
    import { emailOTP } from "better-auth/plugins";
    
    export const auth = betterAuth({
    	plugins: [
    		emailOTP({
    			sendVerificationOTP: async ({ email, otp, type }) => {
    				// Send OTP to user's email
    			}
    		})
    	],
    })

User configuration options.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	user: {
    		modelName: "users",
    		fields: {
    			email: "emailAddress",
    			name: "fullName"
    		},
    		additionalFields: {
    			customField: {
    				type: "string",
    			}
    		},
    		changeEmail: {
    			enabled: true,
    			sendChangeEmailConfirmation: async ({ user, newEmail, url, token }) => {
    				// Send change email confirmation to the old email
    			},
    			updateEmailWithoutVerification: false // Update email without verification if user is not verified
    		},
    		deleteUser: {
    			enabled: true,
    			sendDeleteAccountVerification: async ({ user, url, token }) => {
    				// Send delete account verification
    			},
    			beforeDelete: async (user) => {
    				// Perform actions before user deletion
    			},
    			afterDelete: async (user) => {
    				// Perform cleanup after user deletion
    			}
    		}
    	},
    })

*   `modelName`: The model name for the user (default: `"user"`)
*   `fields`: Map fields to different column names
*   `additionalFields`: Additional fields for the user table
*   `changeEmail`: Configuration for changing email
*   `deleteUser`: Configuration for user deletion

Session configuration options.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	session: {
    		modelName: "sessions",
    		fields: {
    			userId: "user_id"
    		},
    		expiresIn: 604800, // 7 days
    		updateAge: 86400, // 1 day
    		disableSessionRefresh: true, // Disable session refresh so that the session is not updated regardless of the `updateAge` option. (default: `false`)
    		additionalFields: { // Additional fields for the session table
    			customField: {
    				type: "string",
    			}
    		},
    		storeSessionInDatabase: true, // Store session in database when secondary storage is provided (default: `false`)
    		preserveSessionInDatabase: false, // Preserve session records in database when deleted from secondary storage (default: `false`)
    		cookieCache: {
    			enabled: true, // Enable caching session in cookie (default: `false`)	
    			maxAge: 300 // 5 minutes
    		}
    	},
    })

*   `modelName`: The model name for the session (default: `"session"`)
*   `fields`: Map fields to different column names
*   `expiresIn`: Expiration time for the session token in seconds (default: `604800` - 7 days)
*   `updateAge`: How often the session should be refreshed in seconds (default: `86400` - 1 day)
*   `additionalFields`: Additional fields for the session table
*   `storeSessionInDatabase`: Store session in database when secondary storage is provided (default: `false`)
*   `preserveSessionInDatabase`: Preserve session records in database when deleted from secondary storage (default: `false`)
*   `cookieCache`: Enable caching session in cookie

Account configuration options.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	account: {
    		modelName: "accounts",
    		fields: {
    			userId: "user_id"
    		},
    		encryptOAuthTokens: true, // Encrypt OAuth tokens before storing them in the database
    		storeAccountCookie: true, // Store account data after OAuth flow in a cookie (useful for database-less flows)
    		accountLinking: {
    			enabled: true,
    			trustedProviders: ["google", "github", "email-password"],
    			allowDifferentEmails: false
    		}
    	},
    })

*   `modelName`: The model name for the account
*   `fields`: Map fields to different column names

### [`encryptOAuthTokens`](#encryptoauthtokens)

Encrypt OAuth tokens before storing them in the database. Default: `false`.

### [`updateAccountOnSignIn`](#updateaccountonsignin)

If enabled (true), the user account data (accessToken, idToken, refreshToken, etc.) will be updated on sign in with the latest data from the provider.

### [`storeAccountCookie`](#storeaccountcookie)

Store account data after OAuth flow in a cookie. This is useful for database-less flows where you want to store account information (access tokens, refresh tokens, etc.) in a cookie instead of the database.

*   Default: `false`
*   Automatically set to `true` if no database is provided

### [`accountLinking`](#accountlinking)

Configuration for account linking.

*   `enabled`: Enable account linking (default: `false`)
*   `trustedProviders`: List of trusted providers
*   `allowDifferentEmails`: Allow users to link accounts with different email addresses
*   `allowUnlinkingAll`: Allow users to unlink all accounts

Verification configuration options.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	verification: {
    		modelName: "verifications",
    		fields: {
    			userId: "user_id"
    		},
    		disableCleanup: false
    	},
    })

*   `modelName`: The model name for the verification table
*   `fields`: Map fields to different column names
*   `disableCleanup`: Disable cleaning up expired values when a verification value is fetched

Rate limiting configuration.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	rateLimit: {
    		enabled: true,
    		window: 10,
    		max: 100,
    		customRules: {
    			"/example/path": {
    				window: 10,
    				max: 100
    			}
    		},
    		storage: "memory",
    		modelName: "rateLimit"
    	}
    })

*   `enabled`: Enable rate limiting (defaults: `true` in production, `false` in development)
*   `window`: Time window to use for rate limiting. The value should be in seconds. (default: `10`)
*   `max`: The default maximum number of requests allowed within the window. (default: `100`)
*   `customRules`: Custom rate limit rules to apply to specific paths.
*   `storage`: Storage configuration. If you passed a secondary storage, rate limiting will be stored in the secondary storage. (options: `"memory", "database", "secondary-storage"`, default: `"memory"`)
*   `modelName`: The name of the table to use for rate limiting if database is used as storage. (default: `"rateLimit"`)

Advanced configuration options.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	advanced: {
    		ipAddress: {
    			ipAddressHeaders: ["x-client-ip", "x-forwarded-for"],
    			disableIpTracking: false
    		},
    		useSecureCookies: true,
    		disableCSRFCheck: false,
    		crossSubDomainCookies: {
    			enabled: true,
    			additionalCookies: ["custom_cookie"],
    			domain: "example.com"
    		},
    		cookies: {
    			session_token: {
    				name: "custom_session_token",
    				attributes: {
    					httpOnly: true,
    					secure: true
    				}
    			}
    		},
    		defaultCookieAttributes: {
    			httpOnly: true,
    			secure: true
    		},
    		// OAuth state configuration has been moved to account option
    		// Use account.storeStateStrategy and account.skipStateCookieCheck instead
    		cookiePrefix: "myapp",
    		database: {
    			// Use your own custom ID generator,
    			// disable generating IDS so your database will generate them,
    			// or use "serial" to use your database's auto-incrementing ID, or "uuid" to use a random UUID.
    			generateId: (((options: {
    				model: LiteralUnion<Models, string>;
    				size?: number;
    			}) => {
    				return "my-super-unique-id";
    			})) | false | "serial" | "uuid",
    			defaultFindManyLimit: 100,
    			experimentalJoins: false,
    		}
    	},
    })

*   `ipAddress`: IP address configuration for rate limiting and session tracking
*   `useSecureCookies`: Use secure cookies (default: `false`)
*   `disableCSRFCheck`: Disable trusted origins check (⚠️ security risk)
*   `crossSubDomainCookies`: Configure cookies to be shared across subdomains
*   `cookies`: Customize cookie names and attributes
*   `defaultCookieAttributes`: Default attributes for all cookies
*   `cookiePrefix`: Prefix for cookies
*   `database`: Database configuration options
*   OAuth state configuration options (`storeStateStrategy`, `skipStateCookieCheck`) are now part of the `account` option

Logger configuration for Better Auth.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	logger: {
    		disabled: false,
    		disableColors: false,
    		level: "error",
    		log: (level, message, ...args) => {
    			// Custom logging implementation
    			console.log(`[${level}] ${message}`, ...args);
    		}
    	}
    })

The logger configuration allows you to customize how Better Auth handles logging. It supports the following options:

*   `disabled`: Disable all logging when set to `true` (default: `false`)
*   `disableColors`: Disable colors in the default logger implementation (default: determined by the terminal's color support)
*   `level`: Set the minimum log level to display. Available levels are:
    *   `"info"`: Show all logs
    *   `"warn"`: Show warnings and errors
    *   `"error"`: Show only errors
    *   `"debug"`: Show all logs including debug information
*   `log`: Custom logging function that receives:
    *   `level`: The log level (`"info"`, `"warn"`, `"error"`, or `"debug"`)
    *   `message`: The log message
    *   `...args`: Additional arguments passed to the logger

Example with custom logging:

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	logger: {
    		level: "info",
    		log: (level, message, ...args) => {
    			// Send logs to a custom logging service
    			myLoggingService.log({
    				level,
    				message,
    				metadata: args,
    				timestamp: new Date().toISOString()
    			});
    		}
    	}
    })

Database lifecycle hooks for core operations.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	databaseHooks: {
    		user: {
    			create: {
    				before: async (user) => {
    					// Modify user data before creation
    					return { data: { ...user, customField: "value" } };
    				},
    				after: async (user) => {
    					// Perform actions after user creation
    				}
    			},
    			update: {
    				before: async (userData) => {
    					// Modify user data before update
    					return { data: { ...userData, updatedAt: new Date() } };
    				},
    				after: async (user) => {
    					// Perform actions after user update
    				}
    			}
    		},
    		session: {
    			// Session hooks
    		},
    		account: {
    			// Account hooks
    		},
    		verification: {
    			// Verification hooks
    		}
    	},
    })

API error handling configuration.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	onAPIError: {
    		throw: true,
    		onError: (error, ctx) => {
    			// Custom error handling
    			console.error("Auth error:", error);
    		},
    		errorURL: "/auth/error",
    		customizeDefaultErrorPage: {
    			colors: {
    				background: "#ffffff",
    				foreground: "#000000",
    				primary: "#0070f3",
    				primaryForeground: "#ffffff",
    				mutedForeground: "#666666",
    				border: "#e0e0e0",
    				destructive: "#ef4444",
    				titleBorder: "#0070f3",
    				titleColor: "#000000",
    				gridColor: "#f0f0f0",
    				cardBackground: "#ffffff",
    				cornerBorder: "#0070f3"
    			},
    			size: {
    				radiusSm: "0.25rem",
    				radiusMd: "0.5rem",
    				radiusLg: "1rem",
    				textSm: "0.875rem",
    				text2xl: "1.5rem",
    				text4xl: "2.25rem",
    				text6xl: "3.75rem"
    			},
    			font: {
    				defaultFamily: "system-ui, sans-serif",
    				monoFamily: "monospace"
    			},
    			disableTitleBorder: false,
    			disableCornerDecorations: false,
    			disableBackgroundGrid: false
    		}
    	},
    })

*   `throw`: Throw an error on API error (default: `false`)
*   `onError`: Custom error handler
*   `errorURL`: URL to redirect to on error (default: `/api/auth/error`)
*   `customizeDefaultErrorPage`: Configure the default error page provided by Better Auth. Start your dev server and go to `/api/auth/error` to see the error page.
    *   `colors`: Customize color scheme for the error page
        *   `background`: Background color
        *   `foreground`: Foreground/text color
        *   `primary`: Primary accent color
        *   `primaryForeground`: Text color on primary background
        *   `mutedForeground`: Muted text color
        *   `border`: Border color
        *   `destructive`: Error/destructive color
        *   `titleBorder`: Border color for the title
        *   `titleColor`: Title text color
        *   `gridColor`: Background grid color
        *   `cardBackground`: Card background color
        *   `cornerBorder`: Corner decoration border color
    *   `size`: Customize sizing and spacing
        *   `radiusSm`: Small border radius
        *   `radiusMd`: Medium border radius
        *   `radiusLg`: Large border radius
        *   `textSm`: Small text size
        *   `text2xl`: 2xl text size
        *   `text4xl`: 4xl text size
        *   `text6xl`: 6xl text size
    *   `font`: Customize font families
        *   `defaultFamily`: Default font family
        *   `monoFamily`: Monospace font family
    *   `disableTitleBorder`: Disable the border around the title (default: `false`)
    *   `disableCornerDecorations`: Disable corner decorations (default: `false`)
    *   `disableBackgroundGrid`: Disable the background grid pattern (default: `false`)

Request lifecycle hooks.

    import { betterAuth } from "better-auth";
    import { createAuthMiddleware } from "better-auth/api";
    
    export const auth = betterAuth({
    	hooks: {
    		before: createAuthMiddleware(async (ctx) => {
    			// Execute before processing the request
    			console.log("Request path:", ctx.path);
    		}),
    		after: createAuthMiddleware(async (ctx) => {
    			// Execute after processing the request
    			console.log("Response:", ctx.context.returned);
    		})
    	},
    })

For more details and examples, see the [Hooks documentation](https://www.better-auth.com/docs/concepts/hooks).

Disable specific auth paths.

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
    	disabledPaths: ["/sign-up/email", "/sign-in/email"],
    })

Enable or disable Better Auth's telemetry collection. (default: `false`)

    import { betterAuth } from "better-auth";
    export const auth = betterAuth({
      telemetry: {
        enabled: false,
      }
    })</content>
</page>

<page>
  <title>Security | Better Auth</title>
  <url>https://www.better-auth.com/docs/reference/security#trusted-origins</url>
  <content>This page contains information about security features of Better Auth.

Better Auth uses the `scrypt` algorithm to hash passwords by default. This algorithm is designed to be memory-hard and CPU-intensive, making it resistant to brute-force attacks. You can customize the password hashing function by setting the `password` option in the configuration. This option should include a `hash` function to hash passwords and a `verify` function to verify them.

### [Session Expiration](#session-expiration)

Better Auth uses secure session management to protect user data. Sessions are stored in the database or a secondary storage, if configured, to prevent unauthorized access. By default, sessions expire after 7 days, but you can customize this value in the configuration. Additionally, each time a session is used, if it reaches the `updateAge` threshold, the expiration date is extended, which by default is set to 1 day.

### [Session Revocation](#session-revocation)

Better Auth allows you to revoke sessions to enhance security. When a session is revoked, the user is logged out and can no longer access the application. A logged in user can also revoke their own sessions to log out from different devices or browsers.

See the [session management](https://www.better-auth.com/docs/concepts/session-management) for more details.

Better Auth includes multiple safeguards to prevent Cross-Site Request Forgery (CSRF) attacks:

1.  **Avoid simple requests** See [Avoiding simple requests](https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF#avoiding_simple_requests) for more details. Better Auth only allows requests with a non-simple header or a `Content-Type` header of `application/json`.
    
2.  **Origin Validation** Each request’s `Origin` header is verified to confirm it comes from your application or another explicitly trusted source. Requests from untrusted origins are rejected. By default, Better Auth trusts the base URL of your app, but you can specify additional trusted origins via the `trustedOrigins` configuration option.
    
3.  **Secure Cookie Settings** Session cookies use the `SameSite=Lax` attribute by default, preventing browsers from sending cookies with most cross-site requests. You can override this behavior using the `defaultCookieAttributes` option.
    
4.  **No Mutations on GET Requests (with additional safeguards)** `GET` requests are assumed to be read-only and should not alter the application’s state. In cases where a `GET` request must perform a mutation—such as during OAuth callbacks - Better Auth applies extra security measures, including validating `nonce` and `state` parameters to ensure the request’s authenticity.
    

You can skip the CSRF check for all requests by setting the `disableCSRFCheck` option to `true` in the configuration.

    {
      advanced: {
        disableCSRFCheck: true
      }
    }

You can skip the origin check for all requests by setting the `disableOriginCheck` option to `true` in the configuration.

    {
      advanced: {
        disableOriginCheck: true
      }
    }

Skipping csrf check will open your application to CSRF attacks. And skipping origin check may open up your application to other security vulnerabilities including open redirects.

To secure OAuth flows, Better Auth stores the OAuth state and PKCE (Proof Key for Code Exchange) in the database. The state helps prevent CSRF attacks, while PKCE protects against code injection threats. Once the OAuth process completes, these values are removed from the database.

Better Auth assigns secure cookies by default when the base URL uses `https`. These secure cookies are encrypted and only sent over secure connections, adding an extra layer of protection. They are also set with the `sameSite` attribute to `lax` by default to prevent cross-site request forgery attacks. And the `httpOnly` attribute is enabled to prevent client-side JavaScript from accessing the cookie.

For Cross-Subdomain Cookies, you can set the `crossSubDomainCookies` option in the configuration. This option allows cookies to be shared across subdomains, enabling seamless authentication across multiple subdomains.

### [Customizing Cookies](#customizing-cookies)

You can customize cookie names to minimize the risk of fingerprinting attacks and set specific cookie options as needed for additional control. For more information, refer to the [cookie options](https://www.better-auth.com/docs/concepts/cookies).

Plugins can also set custom cookie options to align with specific security needs. If you're using Better Auth in non-browser environments, plugins offer ways to manage cookies securely in those contexts as well.

Better Auth includes built-in rate limiting to safeguard against brute-force attacks. Rate limits are applied across all routes by default, with specific routes subject to stricter limits based on potential risk.

Better Auth uses client IP addresses for rate limiting and security monitoring. By default, it reads the IP address from the standard `X-Forwarded-For` header. However, you can configure a specific trusted header to ensure accurate IP address detection and prevent IP spoofing attacks.

You can configure the IP address header in your Better Auth configuration:

    {
      advanced: {
        ipAddress: {
          ipAddressHeaders: ['cf-connecting-ip'] // or any other custom header
        }
      }
    }

This ensures that Better Auth only accepts IP addresses from your trusted proxy's header, making it more difficult for attackers to bypass rate limiting or other IP-based security measures by spoofing headers.

> **Important**: When setting a custom IP address header, ensure that your proxy or load balancer is properly configured to set this header, and that it cannot be set by end users directly.

Trusted origins prevent CSRF attacks and block open redirects. You can set a list of trusted origins in the `trustedOrigins` configuration option. Requests from origins not on this list are automatically blocked.

### [Basic Usage](#basic-usage)

The most basic usage is to specify exact origins:

    {
      trustedOrigins: [
        "https://example.com",
        "https://app.example.com",
        "http://localhost:3000"
      ]
    }

### [Wildcard Origins](#wildcard-origins)

Better Auth supports wildcard patterns in trusted origins, which allows you to trust multiple subdomains with a single entry:

    {
      trustedOrigins: [
        "*.example.com",             // Trust all subdomains of example.com (any protocol)
        "https://*.example.com",     // Trust only HTTPS subdomains of example.com
        "http://*.dev.example.com"   // Trust all HTTP subdomains of dev.example.com
      ]
    }

#### [Protocol-specific wildcards](#protocol-specific-wildcards)

When using a wildcard pattern with a protocol prefix (like `https://`):

*   The protocol must match exactly
*   The domain can have any subdomain in place of the `*`
*   Requests using a different protocol will be rejected, even if the domain matches

#### [Protocol-agnostic wildcards](#protocol-agnostic-wildcards)

When using a wildcard pattern without a protocol prefix (like `*.example.com`):

*   Any protocol (http, https, etc.) will be accepted
*   The domain must match the wildcard pattern

### [Custom Schemes](#custom-schemes)

Trusted origins also support custom schemes for mobile apps and browser extensions:

    {
      trustedOrigins: [
        "myapp://",                               // Mobile app scheme
        "chrome-extension://YOUR_EXTENSION_ID",   // Browser extension
        "exp://*/*",                              // Trust all Expo development URLs
        "exp://10.0.0.*:*/*",                     // Trust 10.0.0.x IP range with any port
      ]
    }

### [Dynamic origin list](#dynamic-origin-list)

You can also dynamically set the list of trusted origins by providing a function that returns it:

    {
      trustedOrigins: async (request) => {
        const trustedOrigins = await queryTrustedDomains();
        return trustedOrigins;
      }
    }

**Important**: This function will be invoked per incoming request, so be careful if you decide to dynamically fetch your list of trusted domains.

If you discover a security vulnerability in Better Auth, please report it to us at [security@better-auth.com](mailto:security@better-auth.com). We address all reports promptly, and credits will be given for validated discoveries.</content>
</page>

<page>
  <title>Have I Been Pwned | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/have-i-been-pwned</url>
  <content>The Have I Been Pwned plugin helps protect user accounts by preventing the use of passwords that have been exposed in known data breaches. It uses the [Have I Been Pwned](https://haveibeenpwned.com/) API to check if a password has been compromised.

### [Add the plugin to your **auth** config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { haveIBeenPwned } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [
            haveIBeenPwned()
        ]
    })

When a user attempts to create an account or update their password with a compromised password, they'll receive the following default error:

    {
      "code": "PASSWORD_COMPROMISED",
      "message": "Password is compromised"
    }

You can customize the error message:

    haveIBeenPwned({
        customPasswordCompromisedMessage: "Please choose a more secure password."
    })

*   Only the first 5 characters of the password hash are sent to the API
*   The full password is never transmitted
*   Provides an additional layer of account security

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/have-i-been-pwned.mdx)</content>
</page>

<page>
  <title>SolidStart Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/solid-start</url>
  <content>Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Mount the handler](#mount-the-handler)

We need to mount the handler to SolidStart server. Put the following code in your `*auth.ts` file inside `/routes/api/auth` folder.

\*auth.ts

    import { auth } from "~/lib/auth";
    import { toSolidStartHandler } from "better-auth/solid-start";
    
    export const { GET, POST } = toSolidStartHandler(auth);</content>
</page>

<page>
  <title>SAML SSO with Okta | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/saml-sso-with-okta</url>
  <content>This guide walks you through setting up SAML Single Sign-On (SSO) with your Identity Provider (IdP), using Okta as an example. For advanced configuration details and the full API reference, check out the [SSO Plugin Documentation](https://www.better-auth.com/docs/plugins/sso).

SAML (Security Assertion Markup Language) is an XML-based standard for exchanging authentication and authorization data between an Identity Provider (IdP) (e.g., Okta, Azure AD, OneLogin) and a Service Provider (SP) (in this case, Better Auth).

In this setup:

*   **IdP (Okta)**: Authenticates users and sends assertions about their identity.
*   **SP (Better Auth)**: Validates assertions and logs the user in.up.

### [Step 1: Create a SAML Application in Okta](#step-1-create-a-saml-application-in-okta)

1.  Log in to your Okta Admin Console
    
2.  Navigate to Applications > Applications
    
3.  Click "Create App Integration"
    
4.  Select "SAML 2.0" as the Sign-in method
    
5.  Configure the following settings:
    
    *   **Single Sign-on URL**: Your Better Auth ACS endpoint (e.g., `http://localhost:3000/api/auth/sso/saml2/sp/acs/sso`). while `sso` being your providerId
    *   **Audience URI (SP Entity ID)**: Your Better Auth metadata URL (e.g., `http://localhost:3000/api/auth/sso/saml2/sp/metadata`)
    *   **Name ID format**: Email Address or any of your choice.
6.  Download the IdP metadata XML file and certificate
    

### [Step 2: Configure Better Auth](#step-2-configure-better-auth)

Here’s an example configuration for Okta in a dev environment:

    const ssoConfig = {
      defaultSSO: [{
        domain: "localhost:3000", // Your domain
        providerId: "sso",
        samlConfig: {
          // SP Configuration
          issuer: "http://localhost:3000/api/auth/sso/saml2/sp/metadata",
          entryPoint: "https://trial-1076874.okta.com/app/trial-1076874_samltest_1/exktofb0a62hqLAUL697/sso/saml",
          callbackUrl: "/dashboard", // Redirect after successful authentication
          
          // IdP Configuration
          idpMetadata: {
            entityID: "https://trial-1076874.okta.com/app/exktofb0a62hqLAUL697/sso/saml/metadata",
            singleSignOnService: [{
              Binding: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
              Location: "https://trial-1076874.okta.com/app/trial-1076874_samltest_1/exktofb0a62hqLAUL697/sso/saml"
            }],
          },
          cert: `-----BEGIN CERTIFICATE-----
    MIIDqjCCApKgAwIBAgIGAZhVGMeUMA0GCSqGSIb3DQEBCwUAMIGVMQswCQYDVQQGEwJVUzETMBEG
    ...
    [Your Okta Certificate]
    ...
    -----END CERTIFICATE-----`,
          
          // SP Metadata
          spMetadata: {
            metadata: `<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" 
              entityID="http://localhost:3000/api/sso/saml2/sp/metadata">
              ...
              [Your SP Metadata XML]
              ...
            </md:EntityDescriptor>`
          }
        }
      }]
    }

### [Step 3: Multiple Default Providers (Optional)](#step-3-multiple-default-providers-optional)

You can configure multiple SAML providers for different domains:

    const ssoConfig = {
      defaultSSO: [
        {
          domain: "company.com",
          providerId: "company-okta",
          samlConfig: {
            // Okta SAML configuration for company.com
          }
        },
        {
          domain: "partner.com", 
          providerId: "partner-adfs",
          samlConfig: {
            // ADFS SAML configuration for partner.com
          }
        },
        {
          domain: "contractor.org",
          providerId: "contractor-azure",
          samlConfig: {
            // Azure AD SAML configuration for contractor.org
          }
        }
      ]
    }

**Explicit**: Pass providerId directly when signing in. **Domain fallback:** Matches based on the user’s email domain. e.g. [user@company.com](mailto:user@company.com) → matches `company-okta` provider.

### [Step 4: Initiating Sign-In](#step-4-initiating-sign-in)

You can start an SSO flow in three ways:

**1\. Explicitly by `providerId` (recommended):**

    // Explicitly specify which provider to use
    await authClient.signIn.sso({
      providerId: "company-okta",
      callbackURL: "/dashboard"
    });

**2\. By email domain matching:**

    // Automatically matches provider based on email domain
    await authClient.signIn.sso({
      email: "user@company.com",
      callbackURL: "/dashboard"
    });

**3\. By specifying domain:**

    // Explicitly specify domain for matching
    await authClient.signIn.sso({
      domain: "partner.com",
      callbackURL: "/dashboard"
    });

**Important Notes**:

*   DummyIDP should ONLY be used for development and testing
*   Never use these certificates in production
*   The example uses `localhost:3000` - adjust URLs for your environment
*   For production, always use proper IdP providers like Okta, Azure AD, or OneLogin

### [Step 5: Dynamically Registering SAML Providers](#step-5-dynamically-registering-saml-providers)

For dynamic registration, you should register SAML providers using the API. See the [SSO Plugin Documentation](https://www.better-auth.com/docs/plugins/sso#register-a-saml-provider) for detailed registration instructions.

Example registration:

    await authClient.sso.register({
      providerId: "okta-prod",
      issuer: "https://your-domain.com",
      domain: "your-domain.com",
      samlConfig: {
        // Your production SAML configuration
      }
    });

*   [SSO Plugin Documentation](https://www.better-auth.com/docs/plugins/sso)
*   [Okta SAML Documentation](https://developer.okta.com/docs/concepts/saml/)
*   [SAML 2.0 Specification](https://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf)

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/guides/saml-sso-with-okta.mdx)</content>
</page>

<page>
  <title>Sign In With Ethereum (SIWE) | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/siwe</url>
  <content>The Sign in with Ethereum (SIWE) plugin allows users to authenticate using their Ethereum wallets following the [ERC-4361 standard](https://eips.ethereum.org/EIPS/eip-4361). This plugin provides flexibility by allowing you to implement your own message verification and nonce generation logic.

### [Add the Server Plugin](#add-the-server-plugin)

Add the SIWE plugin to your auth configuration:

auth.ts

    import { betterAuth } from "better-auth";
    import { siwe } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [
            siwe({
                domain: "example.com",
                emailDomainName: "example.com", // optional
                anonymous: false, // optional, default is true
                getNonce: async () => {
                    // Implement your nonce generation logic here
                    return "your-secure-random-nonce";
                },
                verifyMessage: async (args) => {
                    // Implement your SIWE message verification logic here
                    // This should verify the signature against the message
                    return true; // return true if signature is valid
                },
                ensLookup: async (args) => {
                    // Optional: Implement ENS lookup for user names and avatars
                    return {
                        name: "user.eth",
                        avatar: "https://example.com/avatar.png"
                    };
                },
            }),
        ],
    });

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the Client Plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { siweClient } from "better-auth/client/plugins";
    
    export const authClient = createAuthClient({
        plugins: [siweClient()],
    });

### [Generate a Nonce](#generate-a-nonce)

Before signing a SIWE message, you need to generate a nonce for the wallet address:

generate-nonce.ts

    const { data, error } = await authClient.siwe.nonce({
      walletAddress: "0x1234567890abcdef1234567890abcdef12345678",
      chainId: 1, // optional for Ethereum mainnet, required for other chains. Defaults to 1
    });
    
    if (data) {
      console.log("Nonce:", data.nonce);
    }

### [Sign In with Ethereum](#sign-in-with-ethereum)

After generating a nonce and creating a SIWE message, verify the signature to authenticate:

sign-in-siwe.ts

    const { data, error } = await authClient.siwe.verify({
      message: "Your SIWE message string",
      signature: "0x...", // The signature from the user's wallet
      walletAddress: "0x1234567890abcdef1234567890abcdef12345678",
      chainId: 1, // optional for Ethereum mainnet, required for other chains. Must match Chain ID in SIWE message
      email: "user@example.com", // optional, required if anonymous is false
    });
    
    if (data) {
      console.log("Authentication successful:", data.user);
    }

### [Chain-Specific Examples](#chain-specific-examples)

Here are examples for different blockchain networks:

ethereum-mainnet.ts

    // Ethereum Mainnet (chainId can be omitted, defaults to 1)
    const { data, error } = await authClient.siwe.verify({
      message,
      signature,
      walletAddress,
      // chainId: 1 (default)
    });

polygon.ts

    // Polygon (chainId REQUIRED)
    const { data, error } = await authClient.siwe.verify({
      message,
      signature,
      walletAddress,
      chainId: 137, // Required for Polygon
    });

arbitrum.ts

    // Arbitrum (chainId REQUIRED)
    const { data, error } = await authClient.siwe.verify({
      message,
      signature,
      walletAddress,
      chainId: 42161, // Required for Arbitrum
    });

base.ts

    // Base (chainId REQUIRED)
    const { data, error } = await authClient.siwe.verify({
      message,
      signature,
      walletAddress,
      chainId: 8453, // Required for Base
    });

The `chainId` must match the Chain ID specified in your SIWE message. Verification will fail with a 401 error if there's a mismatch between the message's Chain ID and the `chainId` parameter.

### [Server Options](#server-options)

The SIWE plugin accepts the following configuration options:

*   **domain**: The domain name of your application (required for SIWE message generation)
*   **emailDomainName**: The email domain name for creating user accounts when not using anonymous mode. Defaults to the domain from your base URL
*   **anonymous**: Whether to allow anonymous sign-ins without requiring an email. Default is `true`
*   **getNonce**: Function to generate a unique nonce for each sign-in attempt. You must implement this function to return a cryptographically secure random string. Must return a `Promise<string>`
*   **verifyMessage**: Function to verify the signed SIWE message. Receives message details and should return `Promise<boolean>`
*   **ensLookup**: Optional function to lookup ENS names and avatars for Ethereum addresses

### [Client Options](#client-options)

The SIWE client plugin doesn't require any configuration options, but you can pass them if needed for future extensibility:

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { siweClient } from "better-auth/client/plugins";
    
    export const authClient = createAuthClient({
      plugins: [
        siweClient({
          // Optional client configuration can go here
        }),
      ],
    });

The SIWE plugin adds a `walletAddress` table to store user wallet associations:

| Field | Type | Description |
| --- | --- | --- |
| id | string | Primary key |
| userId | string | Reference to user.id |
| address | string | Ethereum wallet address |
| chainId | number | Chain ID (e.g., 1 for Ethereum mainnet) |
| isPrimary | boolean | Whether this is the user's primary wallet |
| createdAt | date | Creation timestamp |

Here's a complete example showing how to implement SIWE authentication:

auth.ts

    import { betterAuth } from "better-auth";
    import { siwe } from "better-auth/plugins";
    import { generateRandomString } from "better-auth/crypto";
    import { verifyMessage, createPublicClient, http } from "viem";
    import { mainnet } from "viem/chains";
    
    export const auth = betterAuth({
      database: {
        // your database configuration
      },
      plugins: [
        siwe({
          domain: "myapp.com",
          emailDomainName: "myapp.com",
          anonymous: false,
          getNonce: async () => {
            // Generate a cryptographically secure random nonce
            return generateRandomString(32);
          },
          verifyMessage: async ({ message, signature, address }) => {
            try {
              // Verify the signature using viem (recommended)
              const isValid = await verifyMessage({
                address: address as `0x${string}`,
                message,
                signature: signature as `0x${string}`,
              });
              return isValid;
            } catch (error) {
              console.error("SIWE verification failed:", error);
              return false;
            }
          },
          ensLookup: async ({ walletAddress }) => {
            try {
              // Optional: lookup ENS name and avatar using viem
              // You can use viem's ENS utilities here
              const client = createPublicClient({
                chain: mainnet,
                transport: http(),
              });
    
              const ensName = await client.getEnsName({
                address: walletAddress as `0x${string}`,
              });
    
              const ensAvatar = ensName
                ? await client.getEnsAvatar({
                    name: ensName,
                  })
                : null;
    
              return {
                name: ensName || walletAddress,
                avatar: ensAvatar || "",
              };
            } catch {
              return {
                name: walletAddress,
                avatar: "",
              };
            }
          },
        }),
      ],
    });

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/siwe.mdx)</content>
</page>

<page>
  <title>Open API | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/open-api</url>
  <content>This is a plugin that provides an Open API reference for Better Auth. It shows all endpoints added by plugins and the core. It also provides a way to test the endpoints. It uses [Scalar](https://scalar.com/) to display the Open API reference.

This plugin is still in the early stages of development. We are working on adding more features to it and filling in the gaps.

### [Add the plugin to your **auth** config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { openAPI } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [ 
            openAPI(), 
        ] 
    })

### [Navigate to `/api/auth/reference` to view the Open API reference](#navigate-to-apiauthreference-to-view-the-open-api-reference)

Each plugin endpoints are grouped by the plugin name. The core endpoints are grouped under the `Default` group. And Model schemas are grouped under the `Models` group.

The Open API reference is generated using the [OpenAPI 3.0](https://swagger.io/specification/) specification. You can use the reference to generate client libraries, documentation, and more.

The reference is generated using the [Scalar](https://scalar.com/) library. Scalar provides a way to view and test the endpoints. You can test the endpoints by clicking on the `Try it out` button and providing the required parameters.

### [Generated Schema](#generated-schema)

To get the generated Open API schema directly as JSON, you can do `auth.api.generateOpenAPISchema()`. This will return the Open API schema as a JSON object.

    import { auth } from "~/lib/auth"
    
    const openAPISchema = await auth.api.generateOpenAPISchema()
    console.log(openAPISchema)

### [Using Scalar with Multiple Sources](#using-scalar-with-multiple-sources)

If you're using Scalar for your API documentation, you can add Better Auth as an additional source alongside your main API:

When using Hono with Scalar for OpenAPI documentation, you can integrate Better Auth by adding it as a source:

    app.get("/docs", Scalar({
      pageTitle: "API Documentation", 
      sources: [
        { url: "/api/open-api", title: "API" },
        // Better Auth schema generation endpoint
        { url: "/api/auth/open-api/generate-schema", title: "Auth" },
      ],
    }));

`path` - The path where the Open API reference is served. Default is `/api/auth/reference`. You can change it to any path you like, but keep in mind that it will be appended to the base path of your auth server.

`disableDefaultReference` - If set to `true`, the default Open API reference UI by Scalar will be disabled. Default is `false`.

This allows you to display both your application's API and Better Auth's authentication endpoints in a unified documentation interface.

`theme` - Allows you to change the theme of the OpenAPI reference page. Default is `default`.

`nonce` - Allows you to pass a nonce string to the inline scripts for Content Security Policy (CSP) compliance. Default is `undefined`.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/open-api.mdx)</content>
</page>

<page>
  <title>Astro Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/astro</url>
  <content>Better Auth comes with first class support for Astro. This guide will show you how to integrate Better Auth with Astro.

Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Mount the handler](#mount-the-handler)

To enable Better Auth to handle requests, we need to mount the handler to a catch all API route. Create a file inside `/pages/api/auth` called `[...all].ts` and add the following code:

pages/api/auth/\[...all\].ts

    import { auth } from "~/auth";
    import type { APIRoute } from "astro";
    
    export const ALL: APIRoute = async (ctx) => {
    	// If you want to use rate limiting, make sure to set the 'x-forwarded-for' header to the request headers from the context
    	// ctx.request.headers.set("x-forwarded-for", ctx.clientAddress);
    	return auth.handler(ctx.request);
    };

You can change the path on your better-auth configuration but it's recommended to keep it as `/api/auth/[...all]`

Astro supports multiple frontend frameworks, so you can easily import your client based on the framework you're using.

If you're not using a frontend framework, you can still import the vanilla client.

### [Astro Locals types](#astro-locals-types)

To have types for your Astro locals, you need to set it inside the `env.d.ts` file.

env.d.ts

    
    /// <reference path="../.astro/types.d.ts" />
    
    declare namespace App {
        // Note: 'import {} from ""' syntax does not work in .d.ts files.
        interface Locals {
            user: import("better-auth").User | null;
            session: import("better-auth").Session | null;
        }
    }

### [Middleware](#middleware)

To protect your routes, you can check if the user is authenticated using the `getSession` method in middleware and set the user and session data using the Astro locals with the types we set before. Start by creating a `middleware.ts` file in the root of your project and follow the example below:

middleware.ts

    import { auth } from "@/auth";
    import { defineMiddleware } from "astro:middleware";
    
    export const onRequest = defineMiddleware(async (context, next) => {
        const isAuthed = await auth.api
            .getSession({
                headers: context.request.headers,
            })
    
        if (isAuthed) {
            context.locals.user = isAuthed.user;
            context.locals.session = isAuthed.session;
        } else {
            context.locals.user = null;
            context.locals.session = null;
        }
    
        return next();
    });

### [Getting session on the server inside `.astro` file](#getting-session-on-the-server-inside-astro-file)

You can use `Astro.locals` to check if the user has session and get the user data from the server side. Here is an example of how you can get the session inside an `.astro` file:

    ---
    import { UserCard } from "@/components/user-card";
    
    const session = () => {
        if (Astro.locals.session) {
            return Astro.locals.session;
        } else {
            // Redirect to login page if the user is not authenticated
            return Astro.redirect("/login");
        }
    }
    
    ---
    
    <UserCard initialSession={session} />

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/astro.mdx)</content>
</page>

<page>
  <title>Nitro Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/nitro</url>
  <content>Better Auth can be integrated with your [Nitro Application](https://nitro.build/) (an open source framework to build web servers).

This guide aims to help you integrate Better Auth with your Nitro application in a few simple steps.

Start by scaffolding a new Nitro application using the following command:

Terminal

    npx giget@latest nitro nitro-app --install

This will create the `nitro-app` directory and install all the dependencies. You can now open the `nitro-app` directory in your code editor.

### [Prisma Adapter Setup](#prisma-adapter-setup)

This guide assumes that you have a basic understanding of Prisma. If you are new to Prisma, you can check out the [Prisma documentation](https://www.prisma.io/docs/getting-started).

The `sqlite` database used in this guide will not work in a production environment. You should replace it with a production-ready database like `PostgreSQL`.

For this guide, we will be using the Prisma adapter. You can install prisma client by running the following command:

`prisma` can be installed as a dev dependency using the following command:

Generate a `schema.prisma` file in the `prisma` directory by running the following command:

Terminal

    npx prisma init

You can now replace the contents of the `schema.prisma` file with the following:

prisma/schema.prisma

    generator client {
      provider = "prisma-client-js"
    }
    
    datasource db {
      provider = "sqlite"
      url      = env("DATABASE_URL")
    }
    
    // Will be deleted. Just need it to generate the prisma client
    model Test {
      id   Int    @id @default(autoincrement())
      name String
    }

Ensure that you update the `DATABASE_URL` in your `.env` file to point to the location of your database.

.env

    DATABASE_URL="file:./dev.db"

Run the following command to generate the Prisma client & sync the database:

Terminal

    npx prisma db push

### [Install & Configure Better Auth](#install--configure-better-auth)

Follow steps 1 & 2 from the [installation guide](https://www.better-auth.com/docs/installation) to install Better Auth in your Nitro application & set up the environment variables.

Once that is done, create your Better Auth instance within the `server/utils/auth.ts` file.

server/utils/auth.ts

    import { betterAuth } from "better-auth";
    import { prismaAdapter } from "better-auth/adapters/prisma";
    import { PrismaClient } from "@prisma/client";
    
    const prisma = new PrismaClient();
    export const auth = betterAuth({
      database: prismaAdapter(prisma, { provider: "sqlite" }),
      emailAndPassword: { enabled: true },
    });

### [Update Prisma Schema](#update-prisma-schema)

Use the Better Auth CLI to update your Prisma schema with the required models by running the following command:

Terminal

    npx @better-auth/cli generate --config server/utils/auth.ts

The `--config` flag is used to specify the path to the file where you have created your Better Auth instance.

Head over to the `prisma/schema.prisma` file & save the file to trigger the format on save.

After saving the file, you can run the `npx prisma db push` command to update the database schema.

You can now mount the Better Auth handler in your Nitro application. You can do this by adding the following code to your `server/routes/api/auth/[...all].ts` file:

server/routes/api/auth/\[...all\].ts

    export default defineEventHandler((event) => {
      return auth.handler(toWebRequest(event));
    });

This is a [catch-all](https://nitro.build/guide/routing#catch-all-route) route that will handle all requests to `/api/auth/*`.

### [CORS](#cors)

You can configure CORS for your Nitro app by creating a plugin.

Start by installing the cors package:

You can now create a new file `server/plugins/cors.ts` and add the following code:

server/plugins/cors.ts

    import cors from "cors";
    export default defineNitroPlugin((plugin) => {
      plugin.h3App.use(
        fromNodeMiddleware(
          cors({
            origin: "*",
          }),
        ),
      );
    });

This will enable CORS for all routes. You can customize the `origin` property to allow requests from specific domains. Ensure that the config is in sync with your frontend application.

### [Auth Guard/Middleware](#auth-guardmiddleware)

You can add an auth guard to your Nitro application to protect routes that require authentication. You can do this by creating a new file `server/utils/require-auth.ts` and adding the following code:

server/utils/require-auth.ts

    import { EventHandler, H3Event } from "h3";
    import { fromNodeHeaders } from "better-auth/node";
    
    /**
     * Middleware used to require authentication for a route.
     *
     * Can be extended to check for specific roles or permissions.
     */
    export const requireAuth: EventHandler = async (event: H3Event) => {
      const headers = event.headers;
    
      const session = await auth.api.getSession({
        headers: headers,
      });
      if (!session)
        throw createError({
          statusCode: 401,
          statusMessage: "Unauthorized",
        });
      // You can save the session to the event context for later use
      event.context.auth = session;
    };

You can now use this event handler/middleware in your routes to protect them:

server/routes/api/secret.get.ts

    // Object syntax of the route handler
    export default defineEventHandler({
      // The user has to be logged in to access this route
      onRequest: [requireAuth],
      handler: async (event) => {
        setResponseStatus(event, 201, "Secret data");
        return { message: "Secret data" };
      },
    });

### [Example](#example)

You can find an example of a Nitro application integrated with Better Auth & Prisma [here](https://github.com/BayBreezy/nitrojs-better-auth-prisma).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/nitro.mdx)</content>
</page>

<page>
  <title>TanStack Start Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/tanstack</url>
  <content>This integration guide is assuming you are using TanStack Start.

Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

### [Mount the handler](#mount-the-handler)

We need to mount the handler to a TanStack API endpoint/Server Route. Create a new file: `/src/routes/api/auth/$.ts`

src/routes/api/auth/$.ts

    import { auth } from '@/lib/auth'
    import { createFileRoute } from '@tanstack/react-router'
    
    export const Route = createFileRoute('/api/auth/$')({
      server: {
        handlers: {
          GET: ({ request }) => {
            return auth.handler(request)
          },
          POST: ({ request }) => {
            return auth.handler(request)
          },
        },
      },
    })

### [Usage tips](#usage-tips)

*   We recommend using the client SDK or `authClient` to handle authentication, rather than server actions with `auth.api`.
*   When you call functions that need to set cookies (like `signInEmail` or `signUpEmail`), you'll need to handle cookie setting for TanStack Start. Better Auth provides a `tanstackStartCookies` plugin to automatically handle this for you.

src/lib/auth.ts

    import { betterAuth } from "better-auth";
    import { tanstackStartCookies } from "better-auth/tanstack-start";
    
    export const auth = betterAuth({
        //...your config
        plugins: [tanstackStartCookies()] // make sure this is the last plugin in the array
    })

Now, when you call functions that set cookies, they will be automatically set using TanStack Start's cookie handling system.

    import { auth } from "@/lib/auth"
    
    const signIn = async () => {
        await auth.api.signInEmail({
            body: {
                email: "user@email.com",
                password: "password",
            }
        })
    }

### [Middleware](#middleware)

You can use TanStack Start's middleware to protect routes that require authentication. Create a middleware that checks for a valid session and redirects unauthenticated users to the login page.

src/middleware/auth.ts

    import { redirect } from "@tanstack/react-router";
    import { createMiddleware } from "@tanstack/react-start";
    import { auth } from "./auth";
    
    export const authMiddleware = createMiddleware().server(
        async ({ next, request }) => {
            const session = await auth.api.getSession({ headers: request.headers })
    
            if (!session) {
                throw redirect({ to: "/login" })
            }
    
            return await next()
        }
    );

You can then use this middleware in your route definitions to protect specific routes:

src/routes/dashboard.tsx

    import { createFileRoute } from '@tanstack/react-router'
    import { authMiddleware } from '@/lib/middleware'
    
    export const Route = createFileRoute('/dashboard')({
      component: RouteComponent,
      server: {
        middleware: [authMiddleware],
      },
    })
    
    function RouteComponent() {
      return <div>Hello "/dashboard"!</div>
    }

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/tanstack.mdx)</content>
</page>

<page>
  <title>Single Sign-On (SSO) | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/sso</url>
  <content>`OIDC` `OAuth2` `SSO` `SAML`

Single Sign-On (SSO) allows users to authenticate with multiple applications using a single set of credentials. This plugin supports OpenID Connect (OIDC), OAuth2 providers, and SAML 2.0.

The SAML 2.0 support is in active development and may not be suitable for production use. Please report any issues or bugs on [GitHub](https://github.com/better-auth/better-auth).

### [Install the plugin](#install-the-plugin)

    npm install @better-auth/sso

### [Add Plugin to the server](#add-plugin-to-the-server)

auth.ts

    import { betterAuth } from "better-auth"
    import { sso } from "@better-auth/sso";
    
    const auth = betterAuth({
        plugins: [ 
            sso() 
        ] 
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { ssoClient } from "@better-auth/sso/client"
    
    const authClient = createAuthClient({
        plugins: [ 
            ssoClient() 
        ] 
    })

### [Register an OIDC Provider](#register-an-oidc-provider)

To register an OIDC provider, use the `registerSSOProvider` endpoint and provide the necessary configuration details for the provider.

A redirect URL will be automatically generated using the provider ID. For instance, if the provider ID is `hydra`, the redirect URL would be `{baseURL}/api/auth/sso/callback/hydra`. Note that `/api/auth` may vary depending on your base path configuration.

If you're using Google or other standard OIDC providers, we recommend using Generic OAuth instead. Please check whether Generic OAuth covers your use case before requesting the SSO plugin.

#### [Example](#example)

### [Register a SAML Provider](#register-a-saml-provider)

To register a SAML provider, use the `registerSSOProvider` endpoint with SAML configuration details. The provider will act as a Service Provider (SP) and integrate with your Identity Provider (IdP).

### [Get Service Provider Metadata](#get-service-provider-metadata)

For SAML providers, you can retrieve the Service Provider metadata XML that needs to be configured in your Identity Provider:

get-sp-metadata.ts

    const response = await auth.api.spMetadata({
        query: {
            providerId: "saml-provider",
            format: "xml" // or "json"
        }
    });
    
    const metadataXML = await response.text();
    console.log(metadataXML);

### [Sign In with SSO](#sign-in-with-sso)

To sign in with an SSO provider, you can call `signIn.sso`

You can sign in using the email with domain matching:

sign-in.ts

    const res = await authClient.signIn.sso({
        email: "user@example.com",
        callbackURL: "/dashboard",
    });

or you can specify the domain:

sign-in-domain.ts

    const res = await authClient.signIn.sso({
        domain: "example.com",
        callbackURL: "/dashboard",
    });

You can also sign in using the organization slug if a provider is associated with an organization:

sign-in-org.ts

    const res = await authClient.signIn.sso({
        organizationSlug: "example-org",
        callbackURL: "/dashboard",
    });

Alternatively, you can sign in using the provider's ID:

sign-in-provider-id.ts

    const res = await authClient.signIn.sso({
        providerId: "example-provider-id",
        callbackURL: "/dashboard",
    });

Optionally, you can pass a login hint (for example, an email address or another identifier) to prefill or direct the identity provider:

sign-in-with-login-hint.ts

    const res = await authClient.signIn.sso({
        providerId: "example-provider-id",
        loginHint: "user@example.com",
        callbackURL: "/dashboard",
    });

To use the server API you can use `signInSSO`

sign-in-org.ts

    const res = await auth.api.signInSSO({
        body: {
            organizationSlug: "example-org",
            callbackURL: "/dashboard",
        }
    });

#### [Full method](#full-method)

    const { data, error } = await authClient.signIn.sso({    email: "john@example.com",    organizationSlug: "example-org",    providerId: "example-provider",    domain: "example.com",    callbackURL: "https://example.com/callback", // required    errorCallbackURL: "https://example.com/callback",    newUserCallbackURL: "https://example.com/new-user",    scopes: ["openid", "email", "profile", "offline_access"],    loginHint: "user@example.com",    requestSignUp: true,});

| Prop | Description | Type |
| --- | --- | --- |
| `email?` | 
The email address to sign in with. This is used to identify the issuer to sign in with. It's optional if the issuer is provided.

 | `string` |
| `organizationSlug?` | 

The slug of the organization to sign in with.

 | `string` |
| `providerId?` | 

The ID of the provider to sign in with. This can be provided instead of email or issuer.

 | `string` |
| `domain?` | 

The domain of the provider.

 | `string` |
| `callbackURL` | 

The URL to redirect to after login.

 | `string` |
| `errorCallbackURL?` | 

The URL to redirect to after login.

 | `string` |
| `newUserCallbackURL?` | 

The URL to redirect to after login if the user is new.

 | `string` |
| `scopes?` | 

Scopes to request from the provider.

 | `string[]` |
| `loginHint?` | 

Login hint to send to the identity provider (e.g., email or identifier).

 | `string` |
| `requestSignUp?` | 

Explicitly request sign-up. Useful when disableImplicitSignUp is true for this provider.

 | `boolean` |

Note: If email is provided and loginHint is not specified, email will be sent as the login\_hint to OIDC providers automatically. SAML flows do not support login\_hint.

When a user is authenticated, if the user does not exist, the user will be provisioned using the `provisionUser` function. If the organization provisioning is enabled and a provider is associated with an organization, the user will be added to the organization.

auth.ts

    const auth = betterAuth({
        plugins: [
            sso({
                provisionUser: async (user) => {
                    // provision user
                },
                organizationProvisioning: {
                    disabled: false,
                    defaultRole: "member",
                    getRole: async (user) => {
                        // get role if needed
                    },
                },
            }),
        ],
    });

The SSO plugin provides powerful provisioning capabilities to automatically set up users and manage their organization memberships when they sign in through SSO providers.

### [User Provisioning](#user-provisioning)

User provisioning allows you to run custom logic whenever a user signs in through an SSO provider. This is useful for:

*   Setting up user profiles with additional data from the SSO provider
*   Synchronizing user attributes with external systems
*   Creating user-specific resources
*   Logging SSO sign-ins
*   Updating user information from the SSO provider

auth.ts

    const auth = betterAuth({
        plugins: [
            sso({
                provisionUser: async ({ user, userInfo, token, provider }) => {
                    // Update user profile with SSO data
                    await updateUserProfile(user.id, {
                        department: userInfo.attributes?.department,
                        jobTitle: userInfo.attributes?.jobTitle,
                        manager: userInfo.attributes?.manager,
                        lastSSOLogin: new Date(),
                    });
    
                    // Create user-specific resources
                    await createUserWorkspace(user.id);
    
                    // Sync with external systems
                    await syncUserWithCRM(user.id, userInfo);
    
                    // Log the SSO sign-in
                    await auditLog.create({
                        userId: user.id,
                        action: 'sso_signin',
                        provider: provider.providerId,
                        metadata: {
                            email: userInfo.email,
                            ssoProvider: provider.issuer,
                        },
                    });
                },
            }),
        ],
    });

The `provisionUser` function receives:

*   **user**: The user object from the database
*   **userInfo**: User information from the SSO provider (includes attributes, email, name, etc.)
*   **token**: OAuth2 tokens (for OIDC providers) - may be undefined for SAML
*   **provider**: The SSO provider configuration

### [Organization Provisioning](#organization-provisioning)

Organization provisioning automatically manages user memberships in organizations when SSO providers are linked to specific organizations. This is particularly useful for:

*   Enterprise SSO where each company/domain maps to an organization
*   Automatic role assignment based on SSO attributes
*   Managing team memberships through SSO

#### [Basic Organization Provisioning](#basic-organization-provisioning)

auth.ts

    const auth = betterAuth({
        plugins: [
            sso({
                organizationProvisioning: {
                    disabled: false,           // Enable org provisioning
                    defaultRole: "member",     // Default role for new members
                },
            }),
        ],
    });

#### [Advanced Organization Provisioning with Custom Roles](#advanced-organization-provisioning-with-custom-roles)

auth.ts

    const auth = betterAuth({
        plugins: [
            sso({
                organizationProvisioning: {
                    disabled: false,
                    defaultRole: "member",
                    getRole: async ({ user, userInfo, provider }) => {
                        // Assign roles based on SSO attributes
                        const department = userInfo.attributes?.department;
                        const jobTitle = userInfo.attributes?.jobTitle;
                        
                        // Admins based on job title
                        if (jobTitle?.toLowerCase().includes('manager') || 
                            jobTitle?.toLowerCase().includes('director') ||
                            jobTitle?.toLowerCase().includes('vp')) {
                            return "admin";
                        }
                        
                        // Special roles for IT department
                        if (department?.toLowerCase() === 'it') {
                            return "admin";
                        }
                        
                        // Default to member for everyone else
                        return "member";
                    },
                },
            }),
        ],
    });

#### [Linking SSO Providers to Organizations](#linking-sso-providers-to-organizations)

When registering an SSO provider, you can link it to a specific organization:

register-org-provider.ts

    await auth.api.registerSSOProvider({
        body: {
            providerId: "acme-corp-saml",
            issuer: "https://acme-corp.okta.com",
            domain: "acmecorp.com",
            organizationId: "org_acme_corp_id", // Link to organization
            samlConfig: {
                // SAML configuration...
            },
        },
        headers,
    });

Now when users from `acmecorp.com` sign in through this provider, they'll automatically be added to the "Acme Corp" organization with the appropriate role.

#### [Multiple Organizations Example](#multiple-organizations-example)

You can set up multiple SSO providers for different organizations:

multi-org-setup.ts

    // Acme Corp SAML provider
    await auth.api.registerSSOProvider({
        body: {
            providerId: "acme-corp",
            issuer: "https://acme.okta.com",
            domain: "acmecorp.com",
            organizationId: "org_acme_id",
            samlConfig: { /* ... */ },
        },
        headers,
    });
    
    // TechStart OIDC provider
    await auth.api.registerSSOProvider({
        body: {
            providerId: "techstart-google",
            issuer: "https://accounts.google.com",
            domain: "techstart.io",
            organizationId: "org_techstart_id",
            oidcConfig: { /* ... */ },
        },
        headers,
    });

#### [Organization Provisioning Flow](#organization-provisioning-flow)

1.  **User signs in** through an SSO provider linked to an organization
2.  **User is authenticated** and either found or created in the database
3.  **Organization membership is checked** - if the user isn't already a member of the linked organization
4.  **Role is determined** using either the `defaultRole` or `getRole` function
5.  **User is added** to the organization with the determined role
6.  **User provisioning runs** (if configured) for additional setup

### [Provisioning Best Practices](#provisioning-best-practices)

#### [1\. Idempotent Operations](#1-idempotent-operations)

Make sure your provisioning functions can be safely run multiple times:

    provisionUser: async ({ user, userInfo }) => {
        // Check if already provisioned
        const existingProfile = await getUserProfile(user.id);
        if (!existingProfile.ssoProvisioned) {
            await createUserResources(user.id);
            await markAsProvisioned(user.id);
        }
        
        // Always update attributes (they might change)
        await updateUserAttributes(user.id, userInfo.attributes);
    },

#### [2\. Error Handling](#2-error-handling)

Handle errors gracefully to avoid blocking user sign-in:

    provisionUser: async ({ user, userInfo }) => {
        try {
            await syncWithExternalSystem(user, userInfo);
        } catch (error) {
            // Log error but don't throw - user can still sign in
            console.error('Failed to sync user with external system:', error);
            await logProvisioningError(user.id, error);
        }
    },

#### [3\. Conditional Provisioning](#3-conditional-provisioning)

Only run certain provisioning steps when needed:

    organizationProvisioning: {
        disabled: false,
        getRole: async ({ user, userInfo, provider }) => {
            // Only process role assignment for certain providers
            if (provider.providerId.includes('enterprise')) {
                return determineEnterpriseRole(userInfo);
            }
            return "member";
        },
    },

### [Default SSO Provider](#default-sso-provider)

auth.ts

    const auth = betterAuth({
        plugins: [
            sso({
                defaultSSO: {
                    providerId: "default-saml", // Provider ID for the default provider
                    samlConfig: {
                        issuer: "https://your-app.com",
                        entryPoint: "https://idp.example.com/sso",
                        cert: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
                        callbackUrl: "http://localhost:3000/api/auth/sso/saml2/sp/acs",
                        spMetadata: {
                            entityID: "http://localhost:3000/api/auth/sso/saml2/sp/metadata",
                            metadata: "<!-- Your SP Metadata XML -->",
                        }
                    }
                }
            })
        ]
    });

The defaultSSO provider will be used when:

1.  No matching provider is found in the database

This allows you to test SAML authentication without setting up providers in the database. The defaultSSO provider supports all the same configuration options as regular SAML providers.

### [Service Provider Configuration](#service-provider-configuration)

When registering a SAML provider, you need to provide Service Provider (SP) metadata configuration:

*   **metadata**: XML metadata for the Service Provider
*   **binding**: The binding method, typically "post" or "redirect"
*   **privateKey**: Private key for signing (optional)
*   **privateKeyPass**: Password for the private key (if encrypted)
*   **isAssertionEncrypted**: Whether assertions should be encrypted
*   **encPrivateKey**: Private key for decryption (if encryption is enabled)
*   **encPrivateKeyPass**: Password for the encryption private key

### [Identity Provider Configuration](#identity-provider-configuration)

You also need to provide Identity Provider (IdP) configuration:

*   **metadata**: XML metadata from your Identity Provider
*   **privateKey**: Private key for the IdP communication (optional)
*   **privateKeyPass**: Password for the IdP private key (if encrypted)
*   **isAssertionEncrypted**: Whether assertions from IdP are encrypted
*   **encPrivateKey**: Private key for IdP assertion decryption
*   **encPrivateKeyPass**: Password for the IdP decryption key

### [SAML Attribute Mapping](#saml-attribute-mapping)

Configure how SAML attributes map to user fields:

    mapping: {
        id: "nameID",           // Default: "nameID"
        email: "email",         // Default: "email" or "nameID"
        name: "displayName",    // Default: "displayName"
        firstName: "givenName", // Default: "givenName"
        lastName: "surname",    // Default: "surname"
        extraFields: {
            department: "department",
            role: "jobTitle",
            phone: "telephoneNumber"
        }
    }

[Domain verification](#domain-verification)
-------------------------------------------

Domain verification allows your application to automatically trust a new SSO provider by automatically validating ownership via the associated domain:

Once enabled, make sure you migrate the database schema (again).

See the [Schema](#schema-for-domain-verification) section to add the fields manually.

### [Verify your domain](#verify-your-domain)

When domain verification is enabled, every new SSO provider will be untrusted at first. This means that new sign-ups or sign-ins will be allowed until the domain ownership has been verified.

To verify your ownership over a domain, follow these steps:

#### [Acquire verification token](#acquire-verification-token)

When an SSO provider is registered, a **verification token** will be issued to the provider (it will be returned as part of the response). You can use this token to prove ownership over the domain.

#### [Create `TXT` DNS record](#create-txt-dns-record)

To do this, you'll need to add a `TXT` record to your domain's DNS settings:

*   **Host:** `better-auth-token-{your-provider-id}` (**Note:** This assumes the default token prefix, which can be customized through the `domainVerification.tokenPrefix` option)
*   **Value:** The verification token you were given.

**Save the record and wait for it to propagate.** This can take up to 48 hours, but it's usually much faster.

#### [Submit a validation request](#submit-a-validation-request)

**Once the DNS record has propagated**, you can submit a validation request (See below)

### [Domain validation request](#domain-validation-request)

Once you have configured your domain, you can use your `auth` instance to submit a validation request. This request will either result in a rejection (could not prove your ownership over the domain) or if the verification is successful, your SSO provider domain will be marked as verified.

    const { data, error } = await authClient.sso.verifyDomain({    providerId: "acme-corp", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `providerId` | 
The provider id

 | `string` |

### [Creating a new verification token](#creating-a-new-verification-token)

Every domain verification token will have a default expiry of 1 week since the moment it was issued or the moment when the SSO provider was registered.

After that time, the token will expire and cannot longer be used. When that happens, you can create a new verification token:

POST

/sso/request-domain-verification

    const { data, error } = await authClient.sso.requestDomainVerification({    providerId: "acme-corp", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `providerId` | 
The provider id

 | `string` |

### [SAML Endpoints](#saml-endpoints)

The plugin automatically creates the following SAML endpoints:

*   **SP Metadata**: `/api/auth/sso/saml2/sp/metadata?providerId={providerId}`
*   **SAML Callback**: `/api/auth/sso/saml2/callback/{providerId}`

The plugin requires additional fields in the `ssoProvider` table to store the provider's configuration.

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | A database identifier |
| issuer | string | \- | The issuer identifier |
| domain | string | \- | The domain of the provider |
| oidcConfig | string | \- | The OIDC configuration (JSON string) |
| samlConfig | string | \- | The SAML configuration (JSON string) |
| userId | string | \- | The user ID |
| providerId | string | \- | The provider ID. Used to identify a provider and to generate a redirect URL. |
| organizationId | string | \- | The organization Id. If provider is linked to an organization. |

### [If you have enabled domain verification:](#if-you-have-enabled-domain-verification)

The `ssoProvider` schema is extended as follows:

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| domainVerified | boolean | \- | A flag indicating whether the provider domain has been verified. |

For a detailed guide on setting up SAML SSO with examples for Okta and testing with DummyIDP, see our [SAML SSO with Okta](https://www.better-auth.com/docs/guides/saml-sso-with-okta).

### [Server](#server)

**provisionUser**: A custom function to provision a user when they sign in with an SSO provider.

**organizationProvisioning**: Options for provisioning users to an organization.

**defaultOverrideUserInfo**: Override user info with the provider info by default.

**disableImplicitSignUp**: Disable implicit sign up for new users.

**trustEmailVerified** — Trusts the `email_verified` flag from the provider. ⚠️ Use this with caution — it can lead to account takeover if misused. Only enable this if you know what you are doing or in a controlled environment.

If you want to allow account linking for specific trusted providers, enable the `accountLinking` option in your auth config and specify those providers in the `trustedProviders` list.</content>
</page>

<page>
  <title>Contributing to BetterAuth | Better Auth</title>
  <url>https://www.better-auth.com/docs/reference/contributing</url>
  <content>Thank you for your interest in contributing to Better Auth! This guide is a concise guide to contributing to Better Auth.

[Getting Started](#getting-started)
-----------------------------------

Before diving in, here are a few important resources:

*   Take a look at our existing [issues](https://github.com/better-auth/better-auth/issues) and [pull requests](https://github.com/better-auth/better-auth/pulls)
*   Join our community discussions in [Discord](https://discord.gg/better-auth)

[Development Setup](#development-setup)
---------------------------------------

To get started with development:

Make sure you have [Node.JS](https://nodejs.org/en/download) installed, preferably on LTS.

### [2\. Clone your fork](#2-clone-your-fork)

    # Replace YOUR-USERNAME with your GitHub username
    git clone https://github.com/YOUR-USERNAME/better-auth.git
    cd better-auth

### [3\. Install dependencies](#3-install-dependencies)

Make sure you have [pnpm](https://pnpm.io/installation) installed!

    pnpm install

### [4\. Prepare ENV files](#4-prepare-env-files)

Copy the example env file to create your new `.env` file.

    cp -n ./docs/.env.example ./docs/.env

[Making changes](#making-changes)
---------------------------------

Once you have an idea of what you want to contribute, you can start making changes. Here are some steps to get started:

### [1\. Create a new branch](#1-create-a-new-branch)

    # Make sure you're on main
    git checkout main
    
    # Pull latest changes
    git pull upstream main
    
    # Create and switch to a new branch
    git checkout -b feature/your-feature-name

### [2\. Start development server](#2-start-development-server)

Start the development server:

    pnpm dev

To start the docs server:

    pnpm -F docs dev

### [3\. Make Your Changes](#3-make-your-changes)

*   Make your changes to the codebase.
    
*   Write tests if needed. (Read more about testing [here](https://www.better-auth.com/docs/reference/contributing#testing))
    
*   Update documentation. (Read more about documenting [here](https://www.better-auth.com/docs/reference/contributing#documentation))
    

### [Issues and Bug Fixes](#issues-and-bug-fixes)

*   Check our [GitHub issues](https://github.com/better-auth/better-auth/issues) for tasks labeled `good first issue`
*   When reporting bugs, include steps to reproduce and expected behavior
*   Comment on issues you'd like to work on to avoid duplicate efforts

### [Framework Integrations](#framework-integrations)

We welcome contributions to support more frameworks:

*   Focus on framework-agnostic solutions where possible
*   Keep integrations minimal and maintainable
*   All integrations currently live in the main package

### [Plugin Development](#plugin-development)

*   For core plugins: Open an issue first to discuss your idea
*   For community plugins: Feel free to develop independently
*   Follow our plugin architecture guidelines

### [Documentation](#documentation)

*   Fix typos and errors
*   Add examples and clarify existing content
*   Ensure documentation is up to date with code changes

[Testing](#testing)
-------------------

We use Vitest for testing. Place test files next to the source files they test:

    import { describe, it, expect } from "vitest";
    import { getTestInstance } from "./test-utils/test-instance";
    
    describe("Feature", () => {
        it("should work as expected", async () => {
            const { client } = await getTestInstance();
            // Test code here
            expect(result).toBeDefined();
        });
    });

### [Using the Test Instance Helper](#using-the-test-instance-helper)

The test instance helper now includes improved async context support for managing user sessions:

    const { client, runWithUser, signInWithTestUser } = await getTestInstance();
    
    // Run tests with a specific user context
    await runWithUser("user@example.com", "password", async (headers) => {
        // All client calls within this block will use the user's session
        const response = await client.getSession();
        // headers are automatically applied
    });
    
    // Or use the test user with async context
    const { runWithDefaultUser } = await signInWithTestUser();
    await runWithDefaultUser(async (headers) => {
        // Code here runs with the test user's session context
    });

### [Testing Best Practices](#testing-best-practices)

*   Write clear commit messages
*   Update documentation to reflect your changes
*   Add tests for new features
*   Follow our coding standards
*   Keep pull requests focused on a single change

[Need Help?](#need-help)
------------------------

Don't hesitate to ask for help! You can:

*   Open an [issue](https://github.com/better-auth/better-auth/issues) with questions
*   Join our [community discussions](https://discord.gg/better-auth)
*   Reach out to project maintainers

Thank you for contributing to Better Auth!</content>
</page>

<page>
  <title>Resources | Better Auth</title>
  <url>https://www.better-auth.com/docs/reference/resources</url>
  <content>A curated collection of resources to help you learn and master Better Auth. From blog posts to video tutorials, find everything you need to get started.

[Video tutorials](#video-tutorials)
-----------------------------------

[](https://www.youtube.com/watch?v=lxslnp-ZEMw)

### [The State of Authentication](https://www.youtube.com/watch?v=lxslnp-ZEMw)

**Theo(t3.gg)** explores the current landscape of authentication, discussing trends, challenges, and where the industry is heading.

[](https://www.youtube.com/watch?v=_OApmLmex14)

### [8 Reasons To Try Better Auth](https://www.youtube.com/watch?v=_OApmLmex14)

**CJ** presents 8 compelling reasons why Better Auth is the BEST auth framework he's ever used, demonstrating its superior features and ease of implementation.

reviewshowcaseimplementation

reviewshowcaseimplementation

nextjsimplementationtutorial

[Blog posts](#blog-posts)
-------------------------

organizationsintegrationpayments

multi-tenantzenstackarchitecture</content>
</page>

<page>
  <title>Browser Extension Guide | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/browser-extension-guide</url>
  <content>In this guide, we'll walk you through the steps of creating a browser extension using [Plasmo](https://docs.plasmo.com/) with Better Auth for authentication.

[Setup & Installations](#setup--installations)
----------------------------------------------

Initialize a new Plasmo project with TailwindCSS and a src directory.

    pnpm create plasmo --with-tailwindcss --with-src

Then, install the Better Auth package.

    pnpm add better-auth

To start the Plasmo development server, run the following command.

    pnpm dev

[Configure tsconfig](#configure-tsconfig)
-----------------------------------------

Configure the `tsconfig.json` file to include `strict` mode.

For this demo, we have also changed the import alias from `~` to `@` and set it to the `src` directory.

tsconfig.json

    {
        "compilerOptions": {
            "paths": {
                "@/_": [
                    "./src/_"
                ]
            },
            "strict": true,
            "baseUrl": "."
        }
    }

[Create the client auth instance](#create-the-client-auth-instance)
-------------------------------------------------------------------

Create a new file at `src/auth/auth-client.ts` and add the following code.

auth-client.ts

    import { createAuthClient } from "better-auth/react"
    
    export const authClient = createAuthClient({
        baseURL: "http://localhost:3000" /* Base URL of your Better Auth backend. */,
        plugins: [],
    });

[Configure the manifest](#configure-the-manifest)
-------------------------------------------------

We must ensure the extension knows the URL to the Better Auth backend.

Head to your package.json file, and add the following code.

package.json

    {
        //...
        "manifest": {
            "host_permissions": [
                "https://URL_TO_YOUR_BACKEND" // localhost works too (e.g. http://localhost:3000)
            ]
        }
    }

[You're now ready!](#youre-now-ready)
-------------------------------------

You have now set up Better Auth for your browser extension.

Add your desired UI and create your dream extension!

To learn more about the client Better Auth API, check out the [client documentation](https://www.better-auth.com/docs/concepts/client).

Here's a quick example 😎

src/popup.tsx

    import { authClient } from "./auth/auth-client"
    
    
    function IndexPopup() {
        const {data, isPending, error} = authClient.useSession();
        if(isPending){
            return <>Loading...</>
        }
        if(error){
            return <>Error: {error.message}</>
        }
        if(data){
            return <>Signed in as {data.user.name}</>
        }
    }
    
    export default IndexPopup;

[Bundle your extension](#bundle-your-extension)
-----------------------------------------------

To get a production build, run the following command.

    pnpm build

Head over to [chrome://extensions](chrome://extensions) and enable developer mode.

Click on "Load Unpacked" and navigate to your extension's `build/chrome-mv3-dev` (or `build/chrome-mv3-prod`) directory.

To see your popup, click on the puzzle piece icon on the Chrome toolbar, and click on your extension.

Learn more about [bundling your extension here.](https://docs.plasmo.com/framework#loading-the-extension-in-chrome)

[Configure the server auth instance](#configure-the-server-auth-instance)
-------------------------------------------------------------------------

First, we will need your extension URL.

An extension URL formed like this: `chrome-extension://YOUR_EXTENSION_ID`.

You can find your extension ID at [chrome://extensions](chrome://extensions).

Head to your server's auth file, and make sure that your extension's URL is added to the `trustedOrigins` list.

server.ts

    import { betterAuth } from "better-auth"
    import { auth } from "@/auth/auth"
    
    export const auth = betterAuth({
        trustedOrigins: ["chrome-extension://YOUR_EXTENSION_ID"],
    })

If you're developing multiple extensions or need to support different browser extensions with different IDs, you can use wildcard patterns:

server.ts

    export const auth = betterAuth({
        trustedOrigins: [
            // Support a specific extension ID
            "chrome-extension://YOUR_EXTENSION_ID",
            
            // Or support multiple extensions with wildcard (less secure)
            "chrome-extension://*"
        ],
    })

Using wildcards for extension origins (`chrome-extension://*`) reduces security by trusting all extensions. It's safer to explicitly list each extension ID you trust. Only use wildcards for development and testing.

[That's it!](#thats-it)
-----------------------

Everything is set up! You can now start developing your extension. 🎉

Congratulations! You've successfully created a browser extension using Better Auth and Plasmo. We highly recommend you visit the [Plasmo documentation](https://docs.plasmo.com/) to learn more about the framework.

If you have any questions, feel free to open an issue on our [GitHub repo](https://github.com/better-auth/better-auth/issues), or join our [Discord server](https://discord.gg/better-auth) for support.</content>
</page>

<page>
  <title>JWT | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/jwt</url>
  <content>The JWT plugin provides endpoints to retrieve a JWT token and a JWKS endpoint to verify the token.

This plugin is not meant as a replacement for the session. It's meant to be used for services that require JWT tokens. If you're looking to use JWT tokens for authentication, check out the [Bearer Plugin](https://www.better-auth.com/docs/plugins/bearer).

### [Add the plugin to your **auth** config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { jwt } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [ 
            jwt(), 
        ] 
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

Once you've installed the plugin, you can start using the JWT & JWKS plugin to get the token and the JWKS through their respective endpoints.

### [Retrieve the token](#retrieve-the-token)

There are multiple ways to retrieve JWT tokens:

1.  **Using the client plugin (recommended)**

Add the `jwtClient` plugin to your auth client configuration:

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { jwtClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
      plugins: [
        jwtClient() 
      ]
    })

Then use the client to get JWT tokens:

    const { data, error } = await authClient.token()
    if (error) {
      // handle error
    }
    if (data) {
      const jwtToken = data.token
      // Use this token for authenticated requests to external services
    }

This is the recommended approach for client applications that need JWT tokens for external API authentication.

2.  **Using your session token**

To get the token, call the `/token` endpoint. This will return the following:

      { 
        "token": "ey..."
      }

Make sure to include the token in the `Authorization` header of your requests if the `bearer` plugin is added in your auth configuration.

    await fetch("/api/auth/token", {
      headers: {
        "Authorization": `Bearer ${token}`
      },
    })

3.  **From `set-auth-jwt` header**

When you call `getSession` method, a JWT is returned in the `set-auth-jwt` header, which you can use to send to your services directly.

    await authClient.getSession({
      fetchOptions: {
        onSuccess: (ctx)=>{
          const jwt = ctx.response.headers.get("set-auth-jwt")
        }
      }
    })

### [Verifying the token](#verifying-the-token)

The token can be verified in your own service, without the need for an additional verify call or database check. For this JWKS is used. The public key can be fetched from the `/api/auth/jwks` endpoint.

Since this key is not subject to frequent changes, it can be cached indefinitely. The key ID (`kid`) that was used to sign a JWT is included in the header of the token. In case a JWT with a different `kid` is received, it is recommended to fetch the JWKS again.

      {
        "keys": [
            {
                "crv": "Ed25519",
                "x": "bDHiLTt7u-VIU7rfmcltcFhaHKLVvWFy-_csKZARUEU",
                "kty": "OKP",
                "kid": "c5c7995d-0037-4553-8aee-b5b620b89b23"
            }
        ]
      }

### [OAuth Provider Mode](#oauth-provider-mode)

If you are making your system oAuth compliant (such as when utilizing the OIDC or MCP plugins), you **MUST** disable the `/token` endpoint (oAuth equivalent `/oauth2/token`) and disable setting the jwt header (oAuth equivalent `/oauth2/userinfo`).

auth.ts

    betterAuth({
      disabledPaths: [
        "/token",
      ],
      plugins: [jwt({
        disableSettingJwtHeader: true,
      })]
    })

#### [Example using jose with remote JWKS](#example-using-jose-with-remote-jwks)

    import { jwtVerify, createRemoteJWKSet } from 'jose'
    
    async function validateToken(token: string) {
      try {
        const JWKS = createRemoteJWKSet(
          new URL('http://localhost:3000/api/auth/jwks')
        )
        const { payload } = await jwtVerify(token, JWKS, {
          issuer: 'http://localhost:3000', // Should match your JWT issuer, which is the BASE_URL
          audience: 'http://localhost:3000', // Should match your JWT audience, which is the BASE_URL by default
        })
        return payload
      } catch (error) {
        console.error('Token validation failed:', error)
        throw error
      }
    }
    
    // Usage example
    const token = 'your.jwt.token' // this is the token you get from the /api/auth/token endpoint
    const payload = await validateToken(token)

#### [Example with local JWKS](#example-with-local-jwks)

    import { jwtVerify, createLocalJWKSet } from 'jose'
    
    
    async function validateToken(token: string) {
      try {
        /**
         * This is the JWKS that you get from the /api/auth/
         * jwks endpoint
         */
        const storedJWKS = {
          keys: [{
            //...
          }]
        };
        const JWKS = createLocalJWKSet({
          keys: storedJWKS.data?.keys!,
        })
        const { payload } = await jwtVerify(token, JWKS, {
          issuer: 'http://localhost:3000', // Should match your JWT issuer, which is the BASE_URL
          audience: 'http://localhost:3000', // Should match your JWT audience, which is the BASE_URL by default
        })
        return payload
      } catch (error) {
        console.error('Token validation failed:', error)
        throw error
      }
    }
    
    // Usage example
    const token = 'your.jwt.token' // this is the token you get from the /api/auth/token endpoint
    const payload = await validateToken(token)

### [Remote JWKS Url](#remote-jwks-url)

Disables the `/jwks` endpoint and uses this endpoint in any discovery such as OIDC.

Useful if your JWKS are not managed at `/jwks` or if your jwks are signed with a certificate and placed on your CDN.

NOTE: you **MUST** specify which asymmetric algorithm is used for signing.

auth.ts

    jwt({
      jwks: {
        remoteUrl: "https://example.com/.well-known/jwks.json",
        keyPairConfig: {
          alg: 'ES256',
        },
      }
    })

### [Custom JWKS Path](#custom-jwks-path)

By default, the JWKS endpoint is available at `/jwks`. You can customize this path using the `jwksPath` option.

This is useful when you need to:

*   Follow OAuth 2.0/OIDC conventions (e.g., `/.well-known/jwks.json`)
*   Match existing API conventions in your application
*   Avoid path conflicts with other endpoints

**Server Configuration:**

auth.ts

    jwt({
      jwks: {
        jwksPath: "/.well-known/jwks.json"
      }
    })

**Client Configuration:**

When using a custom `jwksPath` on the server, you **MUST** configure the client with the same path:

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { jwtClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
      plugins: [
        jwtClient({
          jwks: {
            jwksPath: "/.well-known/jwks.json" // Must match server configuration
          }
        })
      ]
    })

Then you can use the `jwks()` method as usual:

    const { data, error } = await authClient.jwks()
    if (data) {
      // Use data.keys to verify JWT tokens
    }

The `jwksPath` configured on the client **MUST** match the server configuration. If they don't match, the client will not be able to fetch the JWKS.

### [Custom Signing](#custom-signing)

This is an advanced feature. Configuration outside of this plugin **MUST** be provided.

Implementers:

*   `remoteUrl` must be defined if using the `sign` function. This shall store all active keys, not just the current one.
*   If using localized approach, ensure server uses the latest private key when rotated. Depending on deployment, the server may need to be restarted.
*   When using remote approach, verify the payload is unchanged after transit. Use integrity validation like CRC32 or SHA256 checks if available.

#### [Localized Signing](#localized-signing)

auth.ts

    jwt({
      jwks: {
        remoteUrl: "https://example.com/.well-known/jwks.json",
        keyPairConfig: {
          alg: 'EdDSA',
        },
      },
      jwt: {
        sign: async (jwtPayload: JWTPayload) => {
          // this is pseudocode
          return await new SignJWT(jwtPayload)
            .setProtectedHeader({
              alg: "EdDSA",
              kid: process.env.currentKid,
              typ: "JWT",
            })
            .sign(process.env.clientPrivateKey);
        },
      },
    })

#### [Remote Signing](#remote-signing)

Useful if you are using a remote Key Management Service such as [Google KMS](https://cloud.google.com/kms/docs/encrypt-decrypt-rsa#kms-encrypt-asymmetric-nodejs), [Amazon KMS](https://docs.aws.amazon.com/kms/latest/APIReference/API_Sign.html), or [Azure Key Vault](https://learn.microsoft.com/en-us/rest/api/keyvault/keys/sign/sign?view=rest-keyvault-keys-7.4&tabs=HTTP).

auth.ts

    jwt({
      jwks: {
        remoteUrl: "https://example.com/.well-known/jwks.json",
        keyPairConfig: {
          alg: 'ES256',
        },
      },
      jwt: {
        sign: async (jwtPayload: JWTPayload) => {
          // this is pseudocode
          const headers = JSON.stringify({ kid: '123', alg: 'ES256', typ: 'JWT' })
          const payload = JSON.stringify(jwtPayload)
          const encodedHeaders = Buffer.from(headers).toString('base64url')
          const encodedPayload = Buffer.from(payload).toString('base64url')
          const hash = createHash('sha256')
          const data = `${encodedHeaders}.${encodedPayload}`
          hash.update(Buffer.from(data))
          const digest = hash.digest()
          const sig = await remoteSign(digest)
          // integrityCheck(sig)
          const jwt = `${data}.${sig}`
          // verifyJwt(jwt)
          return jwt
        },
      },
    })

The JWT plugin adds the following tables to the database:

### [JWKS](#jwks)

Table Name: `jwks`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each web key |
| publicKey | string | \- | The public part of the web key |
| privateKey | string | \- | The private part of the web key |
| createdAt | Date | \- | Timestamp of when the web key was created |
| expiresAt | Date |  | Timestamp of when the web key expires |

You can customize the table name and fields for the `jwks` table. See the [Database concept documentation](https://www.better-auth.com/docs/concepts/database#custom-table-names) for more information on how to customize plugin schema.

### [Algorithm of the Key Pair](#algorithm-of-the-key-pair)

The algorithm used for the generation of the key pair. The default is **EdDSA** with the **Ed25519** curve. Below are the available options:

auth.ts

    jwt({
      jwks: {
        keyPairConfig: {
          alg: "EdDSA",
          crv: "Ed25519"
        }
      }
    })

#### [EdDSA](#eddsa)

*   **Default Curve**: `Ed25519`
*   **Optional Property**: `crv`
    *   Available options: `Ed25519`, `Ed448`
    *   Default: `Ed25519`

#### [ES256](#es256)

*   No additional properties

#### [RSA256](#rsa256)

*   **Optional Property**: `modulusLength`
    *   Expects a number
    *   Default: `2048`

#### [PS256](#ps256)

*   **Optional Property**: `modulusLength`
    *   Expects a number
    *   Default: `2048`

#### [ECDH-ES](#ecdh-es)

*   **Optional Property**: `crv`
    *   Available options: `P-256`, `P-384`, `P-521`
    *   Default: `P-256`

#### [ES512](#es512)

*   No additional properties

### [Disable private key encryption](#disable-private-key-encryption)

By default, the private key is encrypted using AES256 GCM. You can disable this by setting the `disablePrivateKeyEncryption` option to `true`.

For security reasons, it's recommended to keep the private key encrypted.

auth.ts

    jwt({
      jwks: {
        disablePrivateKeyEncryption: true
      }
    })

### [Key Rotation](#key-rotation)

You can enable key rotation by setting the `rotationInterval` option. This will automatically rotate the key pair at the specified interval.

The default value is `undefined` (disabled).

auth.ts

    jwt({
      jwks: {
        rotationInterval: 60 * 60 * 24 * 30, // 30 days
        gracePeriod: 60 * 60 * 24 * 30 // 30 days
      }
    })

*   `rotationInterval`: The interval in seconds to rotate the key pair.
*   `gracePeriod`: The period in seconds to keep the old key pair valid after rotation. This is useful to allow clients to verify tokens signed by the old key pair. The default value is 30 days.

### [Modify JWT payload](#modify-jwt-payload)

By default the entire user object is added to the JWT payload. You can modify the payload by providing a function to the `definePayload` option.

auth.ts

    jwt({
      jwt: {
        definePayload: ({user}) => {
          return {
            id: user.id,
            email: user.email,
            role: user.role
          }
        }
      }
    })

### [Modify Issuer, Audience, Subject or Expiration time](#modify-issuer-audience-subject-or-expiration-time)

If none is given, the `BASE_URL` is used as the issuer and the audience is set to the `BASE_URL`. The expiration time is set to 15 minutes.

auth.ts

    jwt({
      jwt: {
        issuer: "https://example.com",
        audience: "https://example.com",
        expirationTime: "1h",
        getSubject: (session) => {
          // by default the subject is the user id
          return session.user.email
        }
      }
    })

### [Custom Adapter](#custom-adapter)

By default, the JWT plugin stores and retrieves JWKS from your database. You can provide a custom adapter to override this behavior, allowing you to store JWKS in alternative locations such as Redis, external services, or in-memory storage.

auth.ts

    jwt({
      adapter: {
        getJwks: async (ctx) => {
          // Custom implementation to get all JWKS
          // This overrides the default database query
          return await yourCustomStorage.getAllKeys()
        },
        createJwk: async (ctx, webKey) => {
          // Custom implementation to create a new key
          // This overrides the default database insert
          return await yourCustomStorage.createKey(webKey)
        }
      }
    })</content>
</page>

<page>
  <title>Telemetry | Better Auth</title>
  <url>https://www.better-auth.com/docs/reference/telemetry</url>
  <content>Better Auth collects anonymous usage data to help us improve the project. This is optional, transparent, and disabled by default.

Since v1.3.5, Better Auth collects anonymous telemetry data about general usage if enabled.

Telemetry data helps us understand how Better Auth is being used across different environments so we can improve performance, prioritize features, and fix issues more effectively. It guides our decisions on performance optimizations, feature development, and bug fixes. All data is collected completely anonymously and with privacy in mind, and users can opt out at any time. We strive to keep what we collect as transparent as possible.

The following data points may be reported. Everything is anonymous and intended for aggregate insights only.

*   **Anonymous identifier**: A non-reversible hash derived from your project (`package.json` name and optionally `baseURL`). This lets us de‑duplicate events per project without knowing who you are.
*   **Runtime**: `{ name: "node" | "bun" | "deno", version }`.
*   **Environment**: one of `development`, `production`, `test`, or `ci`.
*   **Framework (if detected)**: `{ name, version }` for frameworks like Next.js, Nuxt, Remix, Astro, SvelteKit, etc.
*   **Database (if detected)**: `{ name, version }` for integrations like PostgreSQL, MySQL, SQLite, Prisma, Drizzle, MongoDB, etc.
*   **System info**: platform, OS release, architecture, CPU count/model/speed, total memory, and flags like `isDocker`, `isWSL`, `isTTY`.
*   **Package manager**: `{ name, version }` derived from the npm user agent.
*   **Redacted auth config snapshot**: A minimized, privacy‑preserving view of your `betterAuth` options produced by `getTelemetryAuthConfig`.

We also collect anonymous telemetry from the CLI:

*   **CLI generate (`cli_generate`)**: outcome `generated | overwritten | appended | no_changes | aborted` plus redacted config.
*   **CLI migrate (`cli_migrate`)**: outcome `migrated | no_changes | aborted | unsupported_adapter` plus adapter id (when relevant) and redacted config.

You can audit telemetry locally by setting the `BETTER_AUTH_TELEMETRY_DEBUG=1` environment variable when running your project or by setting `telemetry: { debug: true }` in your auth config. In this debug mode, telemetry events are logged only to the console.

auth.ts

    export const auth = betterAuth({
      telemetry: { 
        debug: true
      } 
    });

All collected data is fully anonymous and only useful in aggregate. It cannot be traced back to any individual source and is accessible only to a small group of core Better Auth maintainers to guide roadmap decisions.

*   **No PII or secrets**: We do not collect emails, usernames, tokens, secrets, client IDs, client secrets, or database URLs.
*   **No full config**: We never send your full `betterAuth` configuration. Instead we send a reduced, redacted snapshot of non‑sensitive toggles and counts.
*   **Redaction by design**: See [detect-auth-config.ts](https://github.com/better-auth/better-auth/blob/main/packages/better-auth/src/telemetry/detectors/detect-auth-config.ts) in the Better Auth source for the exact shape of what is included. It purposely converts sensitive values to booleans, counts, or generic identifiers.

You can enable telemetry collection in your auth config or by setting an environment variable.

*   Via your auth config.
    
    auth.ts
    
        export const auth = betterAuth({
          telemetry: { 
            enabled: true
          } 
        });
    
*   Via an environment variable.
    
    .env
    
        # Enable telemetry
        BETTER_AUTH_TELEMETRY=1
        
        # Disable telemetry
        BETTER_AUTH_TELEMETRY=0
    

### [When is telemetry sent?](#when-is-telemetry-sent)

*   On `betterAuth` initialization (`type: "init"`).
*   On CLI actions: `generate` and `migrate` as described above.

Telemetry is disabled automatically in tests (`NODE_ENV=test`) unless explicitly overridden by internal tooling.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/reference/telemetry.mdx)</content>
</page>

<page>
  <title>NestJS Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/nestjs</url>
  <content>This guide will show you how to integrate Better Auth with [NestJS](https://nestjs.com/).

Before you start, make sure you have a Better Auth instance configured. If you haven't done that yet, check out the [installation](https://www.better-auth.com/docs/installation).

The NestJS integration is **community maintained**. If you encounter any issues, please open them at [nestjs-better-auth](https://github.com/ThallesP/nestjs-better-auth).

[Installation](#installation)
-----------------------------

Install the NestJS integration library:

    npm install @thallesp/nestjs-better-auth

[Basic Setup](#basic-setup)
---------------------------

Currently the library has beta support for Fastify, if you experience any issues with it, please open an issue at [nestjs-better-auth](https://github.com/ThallesP/nestjs-better-auth).

### [1\. Disable Body Parser](#1-disable-body-parser)

Disable NestJS's built-in body parser to allow Better Auth to handle the raw request body:

main.ts

    import { NestFactory } from "@nestjs/core";
    import { AppModule } from "./app.module";
    
    async function bootstrap() {
      const app = await NestFactory.create(AppModule, {
        bodyParser: false, // Required for Better Auth
      });
      await app.listen(process.env.PORT ?? 3000);
    }
    bootstrap();

### [2\. Import AuthModule](#2-import-authmodule)

Import the `AuthModule` in your root module:

app.module.ts

    import { Module } from '@nestjs/common';
    import { AuthModule } from '@thallesp/nestjs-better-auth';
    import { auth } from "./auth"; // Your Better Auth instance
    
    @Module({
      imports: [
        AuthModule.forRoot({ auth }),
      ],
    })
    export class AppModule {}

### [3\. Route Protection](#3-route-protection)

**Global by default**: An `AuthGuard` is registered globally by this module. All routes are protected unless you explicitly allow access.

Use the `Session` decorator to access the user session:

user.controller.ts

    import { Controller, Get } from '@nestjs/common';
    import { Session, UserSession, AllowAnonymous, OptionalAuth } from '@thallesp/nestjs-better-auth';
    
    @Controller('users')
    export class UserController {
      @Get('me')
      async getProfile(@Session() session: UserSession) {
        return { user: session.user };
      }
    
      @Get('public')
      @AllowAnonymous() // Allow anonymous access
      async getPublic() {
        return { message: 'Public route' };
      }
    
      @Get('optional')
      @OptionalAuth() // Authentication is optional
      async getOptional(@Session() session: UserSession) {
        return { authenticated: !!session };
      }
    }

[Full Documentation](#full-documentation)
-----------------------------------------

For comprehensive documentation including decorators, hooks, global guards, and advanced configuration, visit the [NestJS Better Auth repository](https://github.com/thallesp/nestjs-better-auth).</content>
</page>

<page>
  <title>Captcha | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/captcha</url>
  <content>    import { betterAuth } from "better-auth";
    import { captcha } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [ 
            captcha({ 
                provider: "cloudflare-turnstile", // or google-recaptcha, hcaptcha, captchafox
                secretKey: process.env.TURNSTILE_SECRET_KEY!, 
            }), 
        ], 
    });</content>
</page>

<page>
  <title>OIDC Provider | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/oidc-provider</url>
  <content>The **OIDC Provider Plugin** enables you to build and manage your own OpenID Connect (OIDC) provider, granting full control over user authentication without relying on third-party services like Okta or Azure AD. It also allows other services to authenticate users through your OIDC provider.

**Key Features**:

*   **Client Registration**: Register clients to authenticate with your OIDC provider.
*   **Dynamic Client Registration**: Allow clients to register dynamically.
*   **Trusted Clients**: Configure hard-coded trusted clients with optional consent bypass.
*   **Authorization Code Flow**: Support the Authorization Code Flow.
*   **Public Clients**: Support public clients for SPA, mobile apps, CLI tools, etc.
*   **JWKS Endpoint**: Publish a JWKS endpoint to allow clients to verify tokens. (Not fully implemented)
*   **Refresh Tokens**: Issue refresh tokens and handle access token renewal using the `refresh_token` grant.
*   **OAuth Consent**: Implement OAuth consent screens for user authorization, with an option to bypass consent for trusted applications.
*   **UserInfo Endpoint**: Provide a UserInfo endpoint for clients to retrieve user details.

This plugin is in active development and may not be suitable for production use. Please report any issues or bugs on [GitHub](https://github.com/better-auth/better-auth).

### [Mount the Plugin](#mount-the-plugin)

Add the OIDC plugin to your auth config. See [OIDC Configuration](#oidc-configuration) on how to configure the plugin.

auth.ts

    import { betterAuth } from "better-auth";
    import { oidcProvider } from "better-auth/plugins";
    
    const auth = betterAuth({
        plugins: [oidcProvider({
            loginPage: "/sign-in", // path to the login page
            // ...other options
        })]
    })

### [Migrate the Database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the Client Plugin](#add-the-client-plugin)

Add the OIDC client plugin to your auth client config.

    import { createAuthClient } from "better-auth/client";
    import { oidcClient } from "better-auth/client/plugins"
    const authClient = createAuthClient({
        plugins: [oidcClient({
            // Your OIDC configuration
        })]
    })

Once installed, you can utilize the OIDC Provider to manage authentication flows within your application.

### [Register a New Client](#register-a-new-client)

To register a new OIDC client, use the `oauth2.register` method on the client or `auth.api.registerOAuthApplication` on the server.

Notes

By default, client registration requires authentication. Set `allowDynamicClientRegistration: true` to allow public registration. Make sure to add the `oidcClient()` plugin to your auth client configuration.

    const { data, error } = await authClient.oauth2.register({    redirect_uris: ["https://client.example.com/callback"], // required    token_endpoint_auth_method: "client_secret_basic",    grant_types: ["authorization_code"],    response_types: ["code"],    client_name: "My App",    client_uri: "https://client.example.com",    logo_uri: "https://client.example.com/logo.png",    scope: "profile email",    contacts: ["admin@example.com"],    tos_uri: "https://client.example.com/tos",    policy_uri: "https://client.example.com/policy",    jwks_uri: "https://client.example.com/jwks",    jwks: {"keys": [{"kty": "RSA", "alg": "RS256", "use": "sig", "n": "...", "e": "..."}]},    metadata: {"key": "value"},    software_id: "my-software",    software_version: "1.0.0",    software_statement,});

| Prop | Description | Type |
| --- | --- | --- |
| `redirect_uris` | 
A list of redirect URIs.

 | `string[]` |
| `token_endpoint_auth_method?` | 

The authentication method for the token endpoint.

 | `"none" | "client_secret_basic" | "client_secret_post"` |
| `grant_types?` | 

The grant types supported by the application.

 | `("authorization_code" | "implicit" | "password" | "client_credentials" | "refresh_token" | "urn:ietf:params:oauth:grant-type:jwt-bearer" | "urn:ietf:params:oauth:grant-type:saml2-bearer")[]` |
| `response_types?` | 

The response types supported by the application.

 | `("code" | "token")[]` |
| `client_name?` | 

The name of the application.

 | `string` |
| `client_uri?` | 

The URI of the application.

 | `string` |
| `logo_uri?` | 

The URI of the application logo.

 | `string` |
| `scope?` | 

The scopes supported by the application. Separated by spaces.

 | `string` |
| `contacts?` | 

The contact information for the application.

 | `string[]` |
| `tos_uri?` | 

The URI of the application terms of service.

 | `string` |
| `policy_uri?` | 

The URI of the application privacy policy.

 | `string` |
| `jwks_uri?` | 

The URI of the application JWKS.

 | `string` |
| `jwks?` | 

The JWKS of the application.

 | `Record<string, any>` |
| `metadata?` | 

The metadata of the application.

 | `Record<string, any>` |
| `software_id?` | 

The software ID of the application.

 | `string` |
| `software_version?` | 

The software version of the application.

 | `string` |
| `software_statement?` | 

The software statement of the application.

 | `string` |

This endpoint supports [RFC7591](https://datatracker.ietf.org/doc/html/rfc7591) compliant client registration.

Once the application is created, you will receive a `client_id` and `client_secret` that you can display to the user.

### [Trusted Clients](#trusted-clients)

For first-party applications and internal services, you can configure trusted clients directly in your OIDC provider configuration. Trusted clients bypass database lookups for better performance and can optionally skip consent screens for improved user experience.

auth.ts

    import { betterAuth } from "better-auth";
    import { oidcProvider } from "better-auth/plugins";
    
    const auth = betterAuth({
        plugins: [
          oidcProvider({
            loginPage: "/sign-in",
            trustedClients: [
                {
                    clientId: "internal-dashboard",
                    clientSecret: "secure-secret-here",
                    name: "Internal Dashboard",
                    type: "web",
                    redirectURLs: ["https://dashboard.company.com/auth/callback"],
                    disabled: false,
                    skipConsent: true, // Skip consent for this trusted client
                    metadata: { internal: true }
                },
                {
                    clientId: "mobile-app",
                    clientSecret: "mobile-secret", 
                    name: "Company Mobile App",
                    type: "native",
                    redirectURLs: ["com.company.app://auth"],
                    disabled: false,
                    skipConsent: false, // Still require consent if needed
                    metadata: {}
                }
            ]
        })]
    })

### [UserInfo Endpoint](#userinfo-endpoint)

The OIDC Provider includes a UserInfo endpoint that allows clients to retrieve information about the authenticated user. This endpoint is available at `/oauth2/userinfo` and requires a valid access token.

#### [Server-Side Usage](#server-side-usage)

server.ts

    import { auth } from "@/lib/auth";
    
    const userInfo = await auth.api.oAuth2userInfo({
      headers: {
        authorization: "Bearer ACCESS_TOKEN"
      }
    });
    // userInfo contains user details based on the scopes granted

#### [Client-Side Usage (For Third-Party OAuth Clients)](#client-side-usage-for-third-party-oauth-clients)

Third-party OAuth clients can call the UserInfo endpoint using standard HTTP requests:

external-client.ts

    const response = await fetch('https://your-domain.com/api/auth/oauth2/userinfo', {
      headers: {
        'Authorization': 'Bearer ACCESS_TOKEN'
      }
    });
    
    const userInfo = await response.json();

**Returned claims based on scopes:**

*   With `openid` scope: Returns the user's ID (`sub` claim)
*   With `profile` scope: Returns `name`, `picture`, `given_name`, `family_name`
*   With `email` scope: Returns `email` and `email_verified`

#### [Custom Claims](#custom-claims)

The `getAdditionalUserInfoClaim` function receives the user object, requested scopes array, and the client, allowing you to conditionally include claims based on the scopes granted during authorization. These additional claims will be included in both the UserInfo endpoint response and the ID token.

auth.ts

    import { betterAuth } from "better-auth";
    import { oidcProvider } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [
            oidcProvider({
                loginPage: "/sign-in",
                getAdditionalUserInfoClaim: async (user, scopes, client) => {
                    const claims: Record<string, any> = {};
                    
                    // Add custom claims based on scopes
                    if (scopes.includes("profile")) {
                        claims.department = user.department;
                        claims.job_title = user.jobTitle;
                    }
                    
                    // Add claims based on client metadata
                    if (client.metadata?.includeRoles) {
                        claims.roles = user.roles;
                    }
                    
                    return claims;
                }
            })
        ]
    });

### [Consent Screen](#consent-screen)

When a user is redirected to the OIDC provider for authentication, they may be prompted to authorize the application to access their data. This is known as the consent screen. By default, Better Auth will display a sample consent screen. You can customize the consent screen by providing a `consentPage` option during initialization.

**Note**: Trusted clients with `skipConsent: true` will bypass the consent screen entirely, providing a seamless experience for first-party applications.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        plugins: [oidcProvider({
            consentPage: "/path/to/consent/page"
        })]
    })

The plugin will redirect the user to the specified path with `consent_code`, `client_id` and `scope` query parameters. You can use this information to display a custom consent screen. Once the user consents, you can call `oauth2.consent` to complete the authorization.

The consent endpoint supports two methods for passing the consent code:

**Method 1: URL Parameter**

consent-page.ts

    // Get the consent code from the URL
    const params = new URLSearchParams(window.location.search);
    
    // Submit consent with the code in the request body
    const consentCode = params.get('consent_code');
    if (!consentCode) {
    	throw new Error('Consent code not found in URL parameters');
    }
    
    const res = await client.oauth2.consent({
    	accept: true, // or false to deny
    	consent_code: consentCode,
    });

**Method 2: Cookie-Based**

consent-page.ts

    // The consent code is automatically stored in a signed cookie
    // Just submit the consent decision
    const res = await client.oauth2.consent({
    	accept: true, // or false to deny
    	// consent_code not needed when using cookie-based flow
    });

Both methods are fully supported. The URL parameter method works well with mobile apps and third-party contexts, while the cookie-based method provides a simpler implementation for web applications.

### [Handling Login](#handling-login)

When a user is redirected to the OIDC provider for authentication, if they are not already logged in, they will be redirected to the login page. You can customize the login page by providing a `loginPage` option during initialization.

auth.ts

    import { betterAuth } from "better-auth";
    
    export const auth = betterAuth({
        plugins: [oidcProvider({
            loginPage: "/sign-in"
        })]
    })

You don't need to handle anything from your side; when a new session is created, the plugin will handle continuing the authorization flow.

### [OIDC Metadata](#oidc-metadata)

Customize the OIDC metadata by providing a configuration object during initialization.

auth.ts

    import { betterAuth } from "better-auth";
    import { oidcProvider } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [oidcProvider({
            metadata: {
                issuer: "https://your-domain.com",
                authorization_endpoint: "/custom/oauth2/authorize",
                token_endpoint: "/custom/oauth2/token",
                // ...other custom metadata
            }
        })]
    })

### [JWKS Endpoint](#jwks-endpoint)

The OIDC Provider plugin can integrate with the JWT plugin to provide asymmetric key signing for ID tokens verifiable at a JWKS endpoint.

To make your plugin OIDC compliant, you **MUST** disable the `/token` endpoint, the OAuth equivalent is located at `/oauth2/token` instead.

auth.ts

    import { betterAuth } from "better-auth";
    import { oidcProvider } from "better-auth/plugins";
    import { jwt } from "better-auth/plugins";
    
    export const auth = betterAuth({
        disabledPaths: [
            "/token",
        ],
        plugins: [
            jwt(), // Make sure to add the JWT plugin
            oidcProvider({
                useJWTPlugin: true, // Enable JWT plugin integration
                loginPage: "/sign-in",
                // ... other options
            })
        ]
    })

When `useJWTPlugin: false` (default), ID tokens are signed with the application secret.

### [Dynamic Client Registration](#dynamic-client-registration)

If you want to allow clients to register dynamically, you can enable this feature by setting the `allowDynamicClientRegistration` option to `true`.

auth.ts

    const auth = betterAuth({
        plugins: [oidcProvider({
            allowDynamicClientRegistration: true,
        })]
    })

This will allow clients to register using the `/register` endpoint to be publicly available.

The OIDC Provider plugin adds the following tables to the database:

### [OAuth Application](#oauth-application)

Table Name: `oauthApplication`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Database ID of the OAuth client |
| clientId | string |  | Unique identifier for each OAuth client |
| clientSecret | string |  | Secret key for the OAuth client. Optional for public clients using PKCE. |
| name | string | \- | Name of the OAuth client |
| redirectURLs | string | \- | Comma-separated list of redirect URLs |
| metadata | string |  | Additional metadata for the OAuth client |
| type | string | \- | Type of OAuth client (e.g., web, mobile) |
| disabled | boolean | \- | Indicates if the client is disabled |
| userId | string |  | ID of the user who owns the client. (optional) |
| createdAt | Date | \- | Timestamp of when the OAuth client was created |
| updatedAt | Date | \- | Timestamp of when the OAuth client was last updated |

### [OAuth Access Token](#oauth-access-token)

Table Name: `oauthAccessToken`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Database ID of the access token |
| accessToken | string | \- | Access token issued to the client |
| refreshToken | string | \- | Refresh token issued to the client |
| accessTokenExpiresAt | Date | \- | Expiration date of the access token |
| refreshTokenExpiresAt | Date | \- | Expiration date of the refresh token |
| clientId | string |  | ID of the OAuth client |
| userId | string |  | ID of the user associated with the token |
| scopes | string | \- | Comma-separated list of scopes granted |
| createdAt | Date | \- | Timestamp of when the access token was created |
| updatedAt | Date | \- | Timestamp of when the access token was last updated |

### [OAuth Consent](#oauth-consent)

Table Name: `oauthConsent`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Database ID of the consent |
| userId | string |  | ID of the user who gave consent |
| clientId | string |  | ID of the OAuth client |
| scopes | string | \- | Comma-separated list of scopes consented to |
| consentGiven | boolean | \- | Indicates if consent was given |
| createdAt | Date | \- | Timestamp of when the consent was given |
| updatedAt | Date | \- | Timestamp of when the consent was last updated |

**allowDynamicClientRegistration**: `boolean` - Enable or disable dynamic client registration.

**metadata**: `OIDCMetadata` - Customize the OIDC provider metadata.

**loginPage**: `string` - Path to the custom login page.

**consentPage**: `string` - Path to the custom consent page.

**trustedClients**: `(Client & { skipConsent?: boolean })[]` - Array of trusted clients that are configured directly in the provider options. These clients bypass database lookups and can optionally skip consent screens.

**getAdditionalUserInfoClaim**: `(user: User, scopes: string[], client: Client) => Record<string, any>` - Function to get additional user info claims.

**useJWTPlugin**: `boolean` - When `true`, ID tokens are signed using the JWT plugin's asymmetric keys. When `false` (default), ID tokens are signed with HMAC-SHA256 using the application secret.

**schema**: `AuthPluginSchema` - Customize the OIDC provider schema.</content>
</page>

<page>
  <title>System for Cross-domain Identity Management (SCIM) | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/scim</url>
  <content>System for Cross-domain Identity Management ([SCIM](https://simplecloud.info/#Specification)) makes managing identities in multi-domain scenarios easier to support via a standardized protocol. This plugin exposes a [SCIM](https://simplecloud.info/#Specification) server that allows third party identity providers to sync identities to your service.

### [Install the plugin](#install-the-plugin)

    npm install @better-auth/scim

### [Add Plugin to the server](#add-plugin-to-the-server)

auth.ts

    import { betterAuth } from "better-auth"
    import { scim } from "@better-auth/scim"; 
    
    const auth = betterAuth({
        plugins: [ 
            scim() 
        ] 
    })

### [Enable HTTP methods](#enable-http-methods)

SCIM requires the `POST`, `PUT`, `PATCH` and `DELETE` HTTP methods to be supported by your server. For most frameworks, this will work out of the box, but some frameworks may require additional configuration:

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

Upon registration, this plugin will expose compliant [SCIM 2.0](https://simplecloud.info/#Specification) server. Generally, this server is meant to be consumed by a third-party (your identity provider), and will require a:

*   **SCIM base URL**: This should be the fully qualified URL to the SCIM server (e.g `http://your-app.com/api/auth/scim/v2`)
*   **SCIM bearer token**: See [generating a SCIM token](#generating-a-scim-token)

### [Generating a SCIM token](#generating-a-scim-token)

Before your identity provider can start syncing information to your SCIM server, you need to generate a SCIM token that your identity provider will use to authenticate against it.

A SCIM token is a simple bearer token that you can generate:

    const { data, error } = await authClient.scim.generateToken({    providerId: "acme-corp", // required    organizationId: "the-org",});

| Prop | Description | Type |
| --- | --- | --- |
| `providerId` | 
The provider id

 | `string` |
| `organizationId?` | 

Optional organization id. When specified, the organizations plugin must also be enabled

 | `string` |

A `SCIM` token is always restricted to a provider, thus you are required to specify a `providerId`. This can be any provider your instance supports (e.g one of the built-in providers such as `credentials` or an external provider registered through an external plugin such as `@better-auth/sso`). Additionally, when the `organization` plugin is registered, you can optionally restrict the token to an organization via the `organizationId`.

**Important:** By default, any authenticated user with access to your better-auth instance will be able to generate a SCIM token. This can be an important security risk to your application, especially in multi-tenant scenarios. It is highly recommended that you implement [hooks](#hooks) to restrict this access to certain roles or users:

    const userRoles = new Set(["admin"]);
    const userAdminIds = new Set(["some-admin-user-id"]);
    
    scim({
        beforeSCIMTokenGenerated: async ({ user, member, scimToken }) => {
            // IMPORTANT: Use this hook to restrict access to certain roles or users
            // At the very least access must be restricted to admin users (see example below)
    
            const userHasAdmin = member?.role && userRoles.has(member.role);
            const userIsAdmin = userAdminIds.size > 0 && userAdminIds.has(user.id);
    
            if (!userHasAdmin && !userIsAdmin) {
                throw new APIError("FORBIDDEN", { message: "User does not have enough permissions" });
            }
        },
    })

See the [hooks](#hooks) documentation for more details about supported hooks.

### [SCIM endpoints](#scim-endpoints)

The following subset of the specification is currently supported:

#### [List users](#list-users)

Get a list of available users in the database. This is restricted to list only users associated to the same provider and organization than your SCIM token.

Notes

Returns the provisioned SCIM user details. See https://datatracker.ietf.org/doc/html/rfc7644#section-3.4.1

    const data = await auth.api.listSCIMUsers({    query: {        filter: 'userName eq "user-a"',    },    // This endpoint requires a bearer authentication token.    headers: { authorization: 'Bearer <token>' },});

| Prop | Description | Type |
| --- | --- | --- |
| `filter?` | 
SCIM compliant filter expression

 | `string` |

#### [Get user](#get-user)

Get an user from the database. The user will be only returned if it belongs to the same provider and organization than the SCIM token.

GET

/scim/v2/Users/:userId

Notes

Returns the provisioned SCIM user details. See https://datatracker.ietf.org/doc/html/rfc7644#section-3.4.1

    const data = await auth.api.getSCIMUser({    params: {        userId: "user id", // required    },    // This endpoint requires a bearer authentication token.    headers: { authorization: 'Bearer <token>' },});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
Unique user identifier

 | `string` |

#### [Create new user](#create-new-user)

Provisions a new user to the database. The user will have an account associated to the same provider and will be member of the same org than the SCIM token.

Notes

Provision a new user via SCIM. See https://datatracker.ietf.org/doc/html/rfc7644#section-3.3

    const data = await auth.api.createSCIMUser({    body: {        externalId: "third party id",        name: {            formatted: "Daniel Perez",            givenName: "Daniel",            familyName: "Perez",        },        emails: [{ value: "daniel@email.com", primary: true }],    },    // This endpoint requires a bearer authentication token.    headers: { authorization: 'Bearer <token>' },});

| Prop | Description | Type |
| --- | --- | --- |
| `externalId?` | 
Unique external (third party) identifier

 | `string` |
| `name?` | 

User name details

 | `Object` |
| `name.formatted?` | 

Formatted name (takes priority over given and family name)

 | `string` |
| `name.givenName?` | 

Given name

 | `string` |
| `name.familyName?` | 

Family name

 | `string` |
| `emails?` | 

List of emails associated to the user, only a single email can be primary

 | `Array<{ value: string, primary?: boolean }>` |

#### [Update an existing user](#update-an-existing-user)

Replaces an existing user details in the database. This operation can only update users that belong to the same provider and organization than the SCIM token.

PUT

/scim/v2/Users/:userId

Notes

Updates an existing user via SCIM. See https://datatracker.ietf.org/doc/html/rfc7644#section-3.3

    const data = await auth.api.updateSCIMUser({    body: {        externalId: "third party id",        name: {            formatted: "Daniel Perez",            givenName: "Daniel",            familyName: "Perez",        },        emails: [{ value: "daniel@email.com", primary: true }],    },    // This endpoint requires a bearer authentication token.    headers: { authorization: 'Bearer <token>' },});

| Prop | Description | Type |
| --- | --- | --- |
| `externalId?` | 
Unique external (third party) identifier

 | `string` |
| `name?` | 

User name details

 | `Object` |
| `name.formatted?` | 

Formatted name (takes priority over given and family name)

 | `string` |
| `name.givenName?` | 

Given name

 | `string` |
| `name.familyName?` | 

Family name

 | `string` |
| `emails?` | 

List of emails associated to the user, only a single email can be primary

 | `Array<{ value: string, primary?: boolean }>` |

#### [Partial update an existing user](#partial-update-an-existing-user)

Allows to apply a partial update to the user details. This operation can only update users that belong to the same provider and organization than the SCIM token.

PATCH

/scim/v2/Users/:userId

Notes

Partially updates a user resource. See https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2

    const data = await auth.api.patchSCIMUser({    body: {        schemas: ["urn:ietf:params:scim:api:messages:2.0:PatchOp"], // required        Operations: [{ op: "replace", path: "/userName", value: "any value" }], // required    },    // This endpoint requires a bearer authentication token.    headers: { authorization: 'Bearer <token>' },});

| Prop | Description | Type |
| --- | --- | --- |
| `schemas` | 
Mandatory schema declaration

 | `string[]` |
| `Operations` | 

List of JSON patch operations

 | `Array<{ op: "replace" | "add" | "remove", path: string, value: any }>` |

#### [Deletes a user resource](#deletes-a-user-resource)

Completely deletes a user resource from the database. This operation can only delete users that belong to the same provider and organization than the SCIM token.

DELETE

/scim/v2/Users/:userId

Notes

Deletes an existing user resource. See https://datatracker.ietf.org/doc/html/rfc7644#section-3.6

    const data = await auth.api.deleteSCIMUser({    params: {        userId, // required    },    // This endpoint requires a bearer authentication token.    headers: { authorization: 'Bearer <token>' },});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` |  | `string` |

#### [Get service provider config](#get-service-provider-config)

Get SCIM metadata describing supported features of this server.

GET

/scim/v2/ServiceProviderConfig

Notes

Standard SCIM metadata endpoint used by identity providers. See https://datatracker.ietf.org/doc/html/rfc7644#section-4

    const data = await auth.api.getSCIMServiceProviderConfig();

#### [Get SCIM schemas](#get-scim-schemas)

Get the list of supported SCIM schemas.

Notes

Standard SCIM metadata endpoint used by identity providers to acquire information about supported schemas. See https://datatracker.ietf.org/doc/html/rfc7644#section-4

    const data = await auth.api.getSCIMSchemas();

#### [Get SCIM schema](#get-scim-schema)

Get the details of a supported SCIM schema.

GET

/scim/v2/Schemas/:schemaId

Notes

Standard SCIM metadata endpoint used by identity providers to acquire information about a given schema. See https://datatracker.ietf.org/doc/html/rfc7644#section-4

    const data = await auth.api.getSCIMSchema();

#### [Get SCIM resource types](#get-scim-resource-types)

Get the list of supported SCIM types.

GET

/scim/v2/ResourceTypes

Notes

Standard SCIM metadata endpoint used by identity providers to get a list of server supported types. See https://datatracker.ietf.org/doc/html/rfc7644#section-4

    const data = await auth.api.getSCIMResourceTypes();

#### [Get SCIM resource type](#get-scim-resource-type)

Get the details of a supported SCIM resource type.

GET

/scim/v2/ResourceTypes/:resourceTypeId

Notes

Standard SCIM metadata endpoint used by identity providers to get a server supported type. See https://datatracker.ietf.org/doc/html/rfc7644#section-4

    const data = await auth.api.getSCIMResourceType();

#### [SCIM attribute mapping](#scim-attribute-mapping)

By default, the SCIM provisioning will automatically map the following fields:

*   `user.email`: User primary email or the first available email if there is not a primary one
*   `user.name`: Derived from `name` (`name.formatted` or `name.givenName` + `name.familyName`) and fallbacks to the user primary email
*   `account.providerId`: Provider associated to the `SCIM` token
*   `account.accountId`: Defaults to `externalId` and fallbacks to `userName`
*   `member.organizationId`: Organization associated to the provider

The plugin requires additional fields in the `scimProvider` table to store the provider's configuration.

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | A database identifier |
| providerId | string | \- | The provider ID. Used to identify a provider and to generate a redirect URL. |
| scimToken | string | \- | The SCIM bearer token. Used by your identity provider to authenticate against your server |
| organizationId | string | \- | The organization Id. If provider is linked to an organization. |

### [Server](#server)

*   `storeSCIMToken`: The method to store the SCIM token in your database, whether `encrypted`, `hashed` or `plain` text. Default is `plain` text.

Alternatively, you can pass a custom encryptor or hasher to store the SCIM token in your database.

**Custom encryptor**

auth.ts

    scim({
        storeSCIMToken: { 
            encrypt: async (scimToken) => {
                return myCustomEncryptor(scimToken);
            },
            decrypt: async (scimToken) => {
                return myCustomDecryptor(scimToken);
            },
        }
    })

**Custom hasher**

auth.ts

    scim({
        storeSCIMToken: {
            hash: async (scimToken) => {
                return myCustomHasher(scimToken);
            },
        }
    })

### [Hooks](#hooks)

The following hooks allow to intercept the lifecycle of the `SCIM` token generation:

    scim({
        beforeSCIMTokenGenerated: async ({ user, member, scimToken }) => {
            // Callback called before the scim token is persisted
            // can be useful to intercept the generation
            if (member?.role !== "admin") {
                throw new APIError("FORBIDDEN", { message: "User does not have enough permissions" });
            }
        },
        afterSCIMTokenGenerated: async ({ user, member, scimToken, scimProvider }) => {
            // Callback called after the scim token has been persisted
            // can be useful to send a notification or otherwise share the token
            await shareSCIMTokenWithInterestedParty(scimToken);
        },
    })

**Note**: All hooks support error handling. Throwing an error in a before hook will prevent the operation from proceeding.</content>
</page>

<page>
  <title>Create a Database Adapter | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/create-a-db-adapter</url>
  <content>Learn how to create a custom database adapter for Better-Auth using `createAdapterFactory`.

Our `createAdapterFactory` function is designed to be very flexible, and we've done our best to make it easy to understand and use. Our hope is to allow you to focus on writing database logic, and not have to worry about how the adapter is working with Better-Auth.

Anything from custom schema configurations, custom ID generation, safe JSON parsing, key mapping, joins, and more is handled by the `createAdapterFactory` function. All you need to do is provide the database logic, and the `createAdapterFactory` function will handle the rest.

### [Get things ready](#get-things-ready)

1.  Import `createAdapterFactory`.
2.  Create `CustomAdapterConfig` interface that represents your adapter config options.
3.  Create the adapter!

    import { createAdapterFactory, type DBAdapterDebugLogOption } from "better-auth/adapters";
    
    // Your custom adapter config options
    interface CustomAdapterConfig {
      /**
       * Helps you debug issues with the adapter.
       */
      debugLogs?: DBAdapterDebugLogOption;
      /**
       * If the table names in the schema are plural.
       */
      usePlural?: boolean;
    }
    
    export const myAdapter = (config: CustomAdapterConfig = {}) =>
      createAdapterFactory({
        // ...
      });

### [Configure the adapter](#configure-the-adapter)

The `config` object is mostly used to provide information about the adapter to Better-Auth. We try to minimize the amount of code you need to write in your adapter functions, and these `config` options are used to help us do that.

    // ...
    export const myAdapter = (config: CustomAdapterConfig = {}) =>
      createAdapterFactory({
        config: {
          adapterId: "custom-adapter", // A unique identifier for the adapter.
          adapterName: "Custom Adapter", // The name of the adapter.
          usePlural: config.usePlural ?? false, // Whether the table names in the schema are plural.
          debugLogs: config.debugLogs ?? false, // Whether to enable debug logs.
          supportsJSON: false, // Whether the database supports JSON. (Default: false)
          supportsDates: true, // Whether the database supports dates. (Default: true)
          supportsBooleans: true, // Whether the database supports booleans. (Default: true)
          supportsNumericIds: true, // Whether the database supports auto-incrementing numeric IDs. (Default: true)
        },
        // ...
      });

### [Create the adapter](#create-the-adapter)

The `adapter` function is where you write the code that interacts with your database.

    // ...
    export const myAdapter = (config: CustomAdapterConfig = {}) =>
      createAdapterFactory({
        config: {
          // ...
        },
        adapter: ({}) => {
          return {
            create: async ({ data, model, select }) => {
              // ...
            },
            update: async ({ data, model, select }) => {
              // ...
            },
            updateMany: async ({ data, model, select }) => {
              // ...
            },
            delete: async ({ data, model, select }) => {
              // ...
            },
            // ...
          };
        },
      });

Learn more about the `adapter` here [here](https://www.better-auth.com/docs/concepts/database#adapters).

The `adapter` function is where you write the code that interacts with your database.

If you haven't already, check out the `options` object in the [config section](#config), as it can be useful for your adapter.

Before we get into the adapter function, let's go over the parameters that are available to you.

*   `options`: The Better Auth options.
*   `schema`: The schema from the user's Better Auth instance.
*   `debugLog`: The debug log function.
*   `getFieldName`: Function to get the transformed field name for the database.
*   `getModelName`: Function to get the transformed model name for the database.
*   `getDefaultModelName`: Function to get the default model name from the schema.
*   `getDefaultFieldName`: Function to get the default field name from the schema.
*   `getFieldAttributes`: Function to get field attributes for a specific model and field.
*   `transformInput`: Function to transform input data before saving to the database.
*   `transformOutput`: Function to transform output data after retrieving from the database.
*   `transformWhereClause`: Function to transform where clauses for database queries.

Example

    adapter: ({
      options,
      schema,
      debugLog,
      getFieldName,
      getModelName,
      getDefaultModelName,
      getDefaultFieldName,
      getFieldAttributes,
      transformInput,
      transformOutput,
      transformWhereClause,
    }) => {
      return {
        // ...
      };
    };

### [Adapter Methods](#adapter-methods)

*   All `model` values are already transformed into the correct model name for the database based on the end-user's schema configuration.
    *   This also means that if you need access to the `schema` version of a given model, you can't use this exact `model` value, you'll need to use the `getDefaultModelName` function provided in the options to convert the `model` to the `schema` version.
*   We will automatically fill in any missing fields you return based on the user's `schema` configuration.
*   Any method that includes a `select` parameter, is only for the purpose of getting data from your database more efficiently. You do not need to worry about only returning what the `select` parameter states, as we will handle that for you.

### [`create` method](#create-method)

The `create` method is used to create a new record in the database.

It's possible to pass `forceAllowId` as a parameter to the `create` method, which allows `id` to be provided in the `data` object. We handle `forceAllowId` internally, so you don't need to worry about it.

parameters:

*   `model`: The model/table name that new data will be inserted into.
*   `data`: The data to insert into the database.
*   `select`: An array of fields to return from the database.

Make sure to return the data that is inserted into the database.

Example

    create: async ({ model, data, select }) => {
      // Example of inserting data into the database.
      return await db.insert(model).values(data);
    };

### [`update` method](#update-method)

The `update` method is used to update a record in the database.

parameters:

*   `model`: The model/table name that the record will be updated in.
*   `where`: The `where` clause to update the record by.
*   `update`: The data to update the record with.

Make sure to return the data in the row which is updated. This includes any fields that were not updated.

Example

    update: async ({ model, where, update }) => {
      // Example of updating data in the database.
      return await db.update(model).set(update).where(where);
    };

### [`updateMany` method](#updatemany-method)

The `updateMany` method is used to update multiple records in the database.

parameters:

*   `model`: The model/table name that the records will be updated in.
*   `where`: The `where` clause to update the records by.
*   `update`: The data to update the records with.

Make sure to return the number of records that were updated.

Example

    updateMany: async ({ model, where, update }) => {
      // Example of updating multiple records in the database.
      return await db.update(model).set(update).where(where);
    };

### [`delete` method](#delete-method)

The `delete` method is used to delete a record from the database.

parameters:

*   `model`: The model/table name that the record will be deleted from.
*   `where`: The `where` clause to delete the record by.

Example

    delete: async ({ model, where }) => {
      // Example of deleting a record from the database.
      await db.delete(model).where(where);
    }

### [`deleteMany` method](#deletemany-method)

The `deleteMany` method is used to delete multiple records from the database.

parameters:

*   `model`: The model/table name that the records will be deleted from.
*   `where`: The `where` clause to delete the records by.

Make sure to return the number of records that were deleted.

Example

    deleteMany: async ({ model, where }) => {
      // Example of deleting multiple records from the database.
      return await db.delete(model).where(where);
    };

### [`findOne` method](#findone-method)

The `findOne` method is used to find a single record in the database.

parameters:

*   `model`: The model/table name that the record will be found in.
*   `where`: The `where` clause to find the record by.
*   `select`: The `select` clause to return.
*   `join`: Optional join configuration to fetch related records in a single query.

Make sure to return the data that is found in the database.

Example

    findOne: async ({ model, where, select, join }) => {
      // Example of finding a single record in the database.
      return await db.select().from(model).where(where).limit(1);
    };

### [`findMany` method](#findmany-method)

The `findMany` method is used to find multiple records in the database.

parameters:

*   `model`: The model/table name that the records will be found in.
*   `where`: The `where` clause to find the records by.
*   `limit`: The limit of records to return.
*   `sortBy`: The `sortBy` clause to sort the records by.
*   `offset`: The offset of records to return.
*   `join`: Optional join configuration to fetch related records in a single query.

Make sure to return the array of data that is found in the database.

Example

    findMany: async ({ model, where, limit, sortBy, offset, join }) => {
      // Example of finding multiple records in the database.
      return await db
        .select()
        .from(model)
        .where(where)
        .limit(limit)
        .offset(offset)
        .orderBy(sortBy);
    };

### [`count` method](#count-method)

The `count` method is used to count the number of records in the database.

parameters:

*   `model`: The model/table name that the records will be counted in.
*   `where`: The `where` clause to count the records by.

Make sure to return the number of records that were counted.

Example

    count: async ({ model, where }) => {
      // Example of counting the number of records in the database.
      return await db.select().from(model).where(where).count();
    };

### [`options` (optional)](#options-optional)

The `options` object is for any potential config that you got from your custom adapter options.

Example

    const myAdapter = (config: CustomAdapterConfig) =>
      createAdapterFactory({
        config: {
          // ...
        },
        adapter: ({ options }) => {
          return {
            options: config,
          };
        },
      });

### [`createSchema` (optional)](#createschema-optional)

The `createSchema` method allows the [Better Auth CLI](https://www.better-auth.com/docs/concepts/cli) to [generate](https://www.better-auth.com/docs/concepts/cli#generate) a schema for the database.

parameters:

*   `tables`: The tables from the user's Better-Auth instance schema; which is expected to be generated into the schema file.
*   `file`: The file the user may have passed in to the `generate` command as the expected schema file output path.

Example

    createSchema: async ({ file, tables }) => {
      // ... Custom logic to create a schema for the database.
    };

We've provided a test suite that you can use to test your adapter. It requires you to use `vitest`.

my-adapter.test.ts

    import { expect, test, describe } from "vitest";
    import { runAdapterTest } from "better-auth/adapters/test";
    import { myAdapter } from "./my-adapter";
    
    describe("My Adapter Tests", async () => {
      afterAll(async () => {
        // Run DB cleanup here...
      });
      const adapter = myAdapter({
        debugLogs: {
          // If your adapter config allows passing in debug logs, then pass this here.
          isRunningAdapterTests: true, // This is our super secret flag to let us know to only log debug logs if a test fails.
        },
      });
    
      runAdapterTest({
        getAdapter: async (betterAuthOptions = {}) => {
          return adapter(betterAuthOptions);
        },
      });
    });

### [Numeric ID tests](#numeric-id-tests)

If your database supports numeric IDs, then you should run this test as well:

my-adapter.number-id.test.ts

    import { expect, test, describe } from "vitest";
    import { runNumberIdAdapterTest } from "better-auth/adapters/test";
    import { myAdapter } from "./my-adapter";
    
    describe("My Adapter Numeric ID Tests", async () => {
      afterAll(async () => {
        // Run DB cleanup here...
      });
      const adapter = myAdapter({
        debugLogs: {
          // If your adapter config allows passing in debug logs, then pass this here.
          isRunningAdapterTests: true, // This is our super secret flag to let us know to only log debug logs if a test fails.
        },
      });
    
      runNumberIdAdapterTest({
        getAdapter: async (betterAuthOptions = {}) => {
          return adapter(betterAuthOptions);
        },
      });
    });

The `config` object is used to provide information about the adapter to Better-Auth.

We **highly recommend** going through and reading each provided option below, as it will help you understand how to properly configure your adapter.

### [Required Config](#required-config)

### [`adapterId`](#adapterid)

A unique identifier for the adapter.

### [`adapterName`](#adaptername)

The name of the adapter.

### [Optional Config](#optional-config)

### [`supportsNumericIds`](#supportsnumericids)

Whether the database supports numeric IDs. If this is set to `false` and the user's config has enabled `useNumberId`, then we will throw an error.

### [`supportsJSON`](#supportsjson)

Whether the database supports JSON. If the database doesn't support JSON, we will use a `string` to save the JSON data.And when we retrieve the data, we will safely parse the `string` back into a JSON object.

### [`supportsDates`](#supportsdates)

Whether the database supports dates. If the database doesn't support dates, we will use a `string` to save the date. (ISO string) When we retrieve the data, we will safely parse the `string` back into a `Date` object.

### [`supportsBooleans`](#supportsbooleans)

Whether the database supports booleans. If the database doesn't support booleans, we will use a `0` or `1` to save the boolean value. When we retrieve the data, we will safely parse the `0` or `1` back into a boolean value.

### [`usePlural`](#useplural)

Whether the table names in the schema are plural. This is often defined by the user, and passed down through your custom adapter options. If you do not intend to allow the user to customize the table names, you can ignore this option, or set this to `false`.

Example

    const adapter = myAdapter({
      // This value then gets passed into the `usePlural`
      // option in the createAdapterFactory `config` object.
      usePlural: true,
    });

### [`transaction`](#transaction)

Whether the adapter supports transactions. If `false`, operations run sequentially; otherwise provide a function that executes a callback with a `TransactionAdapter`.

If your database does not support transactions, the error handling and rollback will not be as robust. We recommend using a database that supports transactions for better data integrity.

### [`debugLogs`](#debuglogs)

Used to enable debug logs for the adapter. You can pass in a boolean, or an object with the following keys: `create`, `update`, `updateMany`, `findOne`, `findMany`, `delete`, `deleteMany`, `count`. If any of the keys are `true`, the debug logs will be enabled for that method.

Example

    // Will log debug logs for all methods.
    const adapter = myAdapter({
      debugLogs: true,
    });

Example

    // Will only log debug logs for the `create` and `update` methods.
    const adapter = myAdapter({
      debugLogs: {
        create: true,
        update: true,
      },
    });

### [`disableIdGeneration`](#disableidgeneration)

Whether to disable ID generation. If this is set to `true`, then the user's `generateId` option will be ignored.

### [`customIdGenerator`](#customidgenerator)

If your database only supports a specific custom ID generation, then you can use this option to generate your own IDs.

### [`mapKeysTransformInput`](#mapkeystransforminput)

If your database uses a different key name for a given situation, you can use this option to map the keys. This is useful for databases that expect a different key name for a given situation. For example, MongoDB uses `_id` while in Better-Auth we use `id`.

Each key in the returned object represents the old key to replace. The value represents the new key.

This can be a partial object that only transforms some keys.

Example

    mapKeysTransformInput: {
      id: "_id", // We want to replace `id` to `_id` to save into MongoDB
    },

### [`mapKeysTransformOutput`](#mapkeystransformoutput)

If your database uses a different key name for a given situation, you can use this option to map the keys. This is useful for databases that use a different key name for a given situation. For example, MongoDB uses `_id` while in Better-Auth we use `id`.

Each key in the returned object represents the old key to replace. The value represents the new key.

This can be a partial object that only transforms some keys.

Example

    mapKeysTransformOutput: {
      _id: "id", // We want to replace `_id` (from MongoDB) to `id` (for Better-Auth)
    },

### [`customTransformInput`](#customtransforminput)

If you need to transform the input data before it is saved to the database, you can use this option to transform the data.

If you're using `supportsJSON`, `supportsDates`, or `supportsBooleans`, then the transformations will be applied before your `customTransformInput` function is called.

The `customTransformInput` function receives the following arguments:

*   `data`: The data to transform.
*   `field`: The field that is being transformed.
*   `fieldAttributes`: The field attributes of the field that is being transformed.
*   `action`: The action which was called from the adapter (`create` or `update`).
*   `model`: The model that is being transformed.
*   `schema`: The schema that is being transformed.
*   `options`: Better Auth options.

The `customTransformInput` function runs at every key in the data object of a given action.

Example

    customTransformInput: ({ field, data }) => {
      if (field === "id") {
        return "123"; // Force the ID to be "123"
      }
    
      return data;
    };

### [`customTransformOutput`](#customtransformoutput)

If you need to transform the output data before it is returned to the user, you can use this option to transform the data. The `customTransformOutput` function is used to transform the output data. Similar to the `customTransformInput` function, it runs at every key in the data object of a given action, but it runs after the data is retrieved from the database.

Example

    customTransformOutput: ({ field, data }) => {
      if (field === "name") {
        return "Bob"; // Force the name to be "Bob"
      }
    
      return data;
    };

    const some_data = await adapter.create({
      model: "user",
      data: {
        name: "John",
      },
    });
    
    // The name will be "Bob"
    console.log(some_data.name);

### [`disableTransformInput`](#disabletransforminput)

Whether to disable input transformation. This should only be used if you know what you're doing and are manually handling all transformations.

Disabling input transformation can break important adapter functionality like ID generation, boolean/date/JSON conversion, and key mapping.

### [`disableTransformOutput`](#disabletransformoutput)

Whether to disable output transformation. This should only be used if you know what you're doing and are manually handling all transformations.

Disabling output transformation can break important adapter functionality like boolean/date/JSON parsing and key mapping.

### [`disableTransformJoin`](#disabletransformjoin)

Whether to disable join transformation. This should only be used if you know what you're doing and are manually handling joins.

Disabling join transformation can break join functionality.

### [`supportsJoin`](#supportsjoin)

Whether the adapter supports native joins. If set to `false` (the default), Better-Auth will handle joins by making multiple queries and combining the results. If set to `true`, the adapter is expected to handle joins natively (e.g., SQL JOIN operations).

Example

    // For adapters that support native SQL joins
    const adapter = myAdapter({
      supportsJoin: true,
    });</content>
</page>

<page>
  <title>Bearer Token Authentication | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/bearer</url>
  <content>The Bearer plugin enables authentication using Bearer tokens as an alternative to browser cookies. It intercepts requests, adding the Bearer token to the Authorization header before forwarding them to your API.

Use this cautiously; it is intended only for APIs that don't support cookies or require Bearer tokens for authentication. Improper implementation could easily lead to security vulnerabilities.

Add the Bearer plugin to your authentication setup:

auth.ts

    import { betterAuth } from "better-auth";
    import { bearer } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [bearer()]
    });

### [1\. Obtain the Bearer Token](#1-obtain-the-bearer-token)

After a successful sign-in, you'll receive a session token in the response headers. Store this token securely (e.g., in `localStorage`):

auth-client.ts

    const { data } = await authClient.signIn.email({
        email: "user@example.com",
        password: "securepassword"
    }, {
      onSuccess: (ctx)=>{
        const authToken = ctx.response.headers.get("set-auth-token") // get the token from the response headers
        // Store the token securely (e.g., in localStorage)
        localStorage.setItem("bearer_token", authToken);
      }
    });

You can also set this up globally in your auth client:

auth-client.ts

    export const authClient = createAuthClient({
        fetchOptions: {
            onSuccess: (ctx) => {
                const authToken = ctx.response.headers.get("set-auth-token") // get the token from the response headers
                // Store the token securely (e.g., in localStorage)
                if(authToken){
                  localStorage.setItem("bearer_token", authToken);
                }
            }
        }
    });

You may want to clear the token based on the response status code or other conditions:

### [2\. Configure the Auth Client](#2-configure-the-auth-client)

Set up your auth client to include the Bearer token in all requests:

auth-client.ts

    export const authClient = createAuthClient({
        fetchOptions: {
            auth: {
               type:"Bearer",
               token: () => localStorage.getItem("bearer_token") || "" // get the token from localStorage
            }
        }
    });

### [3\. Make Authenticated Requests](#3-make-authenticated-requests)

Now you can make authenticated API calls:

auth-client.ts

    // This request is automatically authenticated
    const { data } = await authClient.listSessions();

### [4\. Per-Request Token (Optional)](#4-per-request-token-optional)

You can also provide the token for individual requests:

auth-client.ts

    const { data } = await authClient.listSessions({
        fetchOptions: {
            headers: {
                Authorization: `Bearer ${token}`
            }
        }
    });

### [5\. Using Bearer Tokens Outside the Auth Client](#5-using-bearer-tokens-outside-the-auth-client)

The Bearer token can be used to authenticate any request to your API, even when not using the auth client:

api-call.ts

    const token = localStorage.getItem("bearer_token");
    
    const response = await fetch("https://api.example.com/data", {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    const data = await response.json();

And in the server, you can use the `auth.api.getSession` function to authenticate requests:

server.ts

    import { auth } from "@/auth";
    
    export async function handler(req, res) {
      const session = await auth.api.getSession({
        headers: req.headers
      });
      
      if (!session) {
        return res.status(401).json({ error: "Unauthorized" });
      }
      
      // Process authenticated request
      // ...
    }

**requireSignature** (boolean): Require the token to be signed. Default: `false`.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/bearer.mdx)</content>
</page>

<page>
  <title>Convex Integration | Better Auth</title>
  <url>https://www.better-auth.com/docs/integrations/convex</url>
  <content>This documentation comes from the [Convex documentation](https://convex-better-auth.netlify.app/), for more information, please refer to their documentation.

### [Create a Convex project](#create-a-convex-project)

To use Convex + Better Auth, you'll first need a [Convex](https://www.convex.dev/) project. If you don't have one, run the following command to get started.

Check out the [Convex docs](https://docs.convex.dev/home) to learn more about Convex.

### [Run `convex dev`](#run-convex-dev)

Running the CLI during setup will initialize your Convex deployment if it doesn't already exist, and keeps generated types current through the process. Keep it running.

The following documentation assumes you're using Next.js.

If you're not using Next.js, support for other frameworks is documented in the [installation guide by Convex](https://convex-better-auth.netlify.app/#select-your-framework).

For a complete example, check out Convex + Better Auth example with Next.js [on GitHub](https://github.com/get-convex/better-auth/tree/main/examples/next).

### [Installation](#installation)

#### [Install packages](#install-packages)

Install the component, a pinned version of Better Auth, and ensure the latest version of Convex.

This component requires Convex `1.25.0` or later.

#### [Register the component](#register-the-component)

Register the Better Auth component in your Convex project.

convex/convex.config.ts

    import { defineApp } from "convex/server";
    import betterAuth from "@convex-dev/better-auth/convex.config";
    
    const app = defineApp();
    app.use(betterAuth);
    
    export default app;

#### [Add Convex auth config](#add-convex-auth-config)

Add a `convex/auth.config.ts` file to configure Better Auth as an authentication provider.

convex/auth.config.ts

    export default {
        providers: [
            {
                domain: process.env.CONVEX_SITE_URL,
                applicationID: "convex",
            },
        ],
    };

#### [Set environment variables](#set-environment-variables)

Generate a secret for encryption and generating hashes. Use the command below if you have openssl installed, or generate your own however you like.

Add your site URL to your Convex deployment.

Add environment variables to the `.env.local` file created by `npx convex dev`. It will be picked up by your framework dev server.

### [Create a Better Auth instance](#create-a-better-auth-instance)

Create a Better Auth instance and initialize the component.

Some TypeScript errors will show until you save the file.

convex/auth.ts

    import { createClient, type GenericCtx } from "@convex-dev/better-auth";
    import { convex } from "@convex-dev/better-auth/plugins";
    import { components } from "./_generated/api";
    import { DataModel } from "./_generated/dataModel";
    import { query } from "./_generated/server";
    import { betterAuth } from "better-auth";
    
    const siteUrl = process.env.SITE_URL!;
    
    // The component client has methods needed for integrating Convex with Better Auth,
    // as well as helper methods for general use.
    export const authComponent = createClient<DataModel>(components.betterAuth);
    
    export const createAuth = (
        ctx: GenericCtx<DataModel>,
        { optionsOnly } = { optionsOnly: false },
    ) => {
        return betterAuth({
            // disable logging when createAuth is called just to generate options.
            // this is not required, but there's a lot of noise in logs without it.
            logger: {
                disabled: optionsOnly,
            },
            baseURL: siteUrl,
            database: authComponent.adapter(ctx),
            // Configure simple, non-verified email/password to get started
            emailAndPassword: {
                enabled: true,
                requireEmailVerification: false,
            },
            plugins: [
                // The Convex plugin is required for Convex compatibility
                convex(),
            ],
        });
    };
    
    // Example function for getting the current user
    // Feel free to edit, omit, etc.
    export const getCurrentUser = query({
        args: {},
        handler: async (ctx) => {
            return authComponent.getAuthUser(ctx);
        },
    });

### [Create a Better Auth client instance](#create-a-better-auth-client-instance)

Create a Better Auth client instance for interacting with the Better Auth server from your client.

src/lib/auth-client.ts

    import { createAuthClient } from "better-auth/react";
    import { convexClient } from "@convex-dev/better-auth/client/plugins";
    
    export const authClient = createAuthClient({
        plugins: [convexClient()],
    });

### [Mount handlers](#mount-handlers)

Register Better Auth route handlers on your Convex deployment.

convex/http.ts

    import { httpRouter } from "convex/server";
    import { authComponent, createAuth } from "./auth";
    
    const http = httpRouter();
    
    authComponent.registerRoutes(http, createAuth);
    
    export default http;

Set up route handlers to proxy auth requests from your framework server to your Convex deployment.

app/api/auth/\[...all\]/route.ts

    import { nextJsHandler } from "@convex-dev/better-auth/nextjs";
    
    export const { GET, POST } = nextJsHandler();

### [Set up Convex client provider](#set-up-convex-client-provider)

Wrap your app with the `ConvexBetterAuthProvider` component.

app/ConvexClientProvider.tsx

    "use client";
    
    import { ReactNode } from "react";
    import { ConvexReactClient } from "convex/react";
    import { authClient } from "@/lib/auth-client"; 
    import { ConvexBetterAuthProvider } from "@convex-dev/better-auth/react"; 
    
    const convex = new ConvexReactClient(process.env.NEXT_PUBLIC_CONVEX_URL!, {
      // Optionally pause queries until the user is authenticated
      expectAuth: true, 
    });
    
    export function ConvexClientProvider({ children }: { children: ReactNode }) {
      return (
        <ConvexBetterAuthProvider client={convex} authClient={authClient}>
          {children}
        </ConvexBetterAuthProvider>
      );
    }

### [You're done!](#youre-done)

You're now ready to start using Better Auth with Convex.

### [Using Better Auth from the server](#using-better-auth-from-the-server)

To use Better Auth's [server methods](https://www.better-auth.com/docs/concepts/api) in server rendering, server functions, or any other Next.js server code, use Convex functions and call the function from your server code.

First, a token helper for calling Convex functions from your server code.

src/lib/auth-server.ts

    import { createAuth } from "@/convex/auth";
    import { getToken as getTokenNextjs } from "@convex-dev/better-auth/nextjs";
    
    export const getToken = () => {
      return getTokenNextjs(createAuth);
    };

Here's an example Convex function that uses Better Auth's server methods, and a server action that calls the Convex function.

convex/users.ts

    import { mutation } from "./_generated/server";
    import { v } from "convex/values";
    import { createAuth, authComponent } from "./auth";
    
    export const updateUserPassword = mutation({
      args: {
        currentPassword: v.string(),
        newPassword: v.string(),
      },
      handler: async (ctx, args) => {
        const { auth, headers } = await authComponent.getAuth(createAuth, ctx);
        await auth.api.changePassword({
          body: {
            currentPassword: args.currentPassword,
            newPassword: args.newPassword,
          },
          headers,
        });
      },
    });

app/actions.ts

    "use server";
    
    import { fetchMutation } from "convex/nextjs";
    import { api } from "../convex/_generated/api";
    import { getToken } from "../lib/auth-server";
    
    // Authenticated mutation via server function
    export async function updatePassword({
      currentPassword,
      newPassword,
    }: {
      currentPassword: string;
      newPassword: string;
    }) {
      const token = await getToken();
      await fetchMutation(
        api.users.updatePassword,
        { currentPassword, newPassword },
        { token }
      );
    }

This documentation comes from the [Convex documentation](https://convex-better-auth.netlify.app/), for more information, please refer to their documentation.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/integrations/convex.mdx)</content>
</page>

<page>
  <title>Device Authorization | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/device-authorization</url>
  <content>`RFC 8628` `CLI` `Smart TV` `IoT`

The Device Authorization plugin implements the OAuth 2.0 Device Authorization Grant ([RFC 8628](https://datatracker.ietf.org/doc/html/rfc8628)), enabling authentication for devices with limited input capabilities such as smart TVs, CLI applications, IoT devices, and gaming consoles.

You can test the device authorization flow right now using the Better Auth CLI:

    npx @better-auth/cli login

This will demonstrate the complete device authorization flow by:

1.  Requesting a device code from the Better Auth demo server
2.  Displaying a user code for you to enter
3.  Opening your browser to the verification page
4.  Polling for authorization completion

The CLI login command is a demo feature that connects to the Better Auth demo server to showcase the device authorization flow in action.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

Add the device authorization plugin to your server configuration.

auth.ts

    import { betterAuth } from "better-auth";
    import { deviceAuthorization } from "better-auth/plugins"; 
    
    export const auth = betterAuth({
      // ... other config
      plugins: [
        deviceAuthorization({ 
          verificationUri: "/device", 
        }), 
      ],
    });

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

Add the device authorization plugin to your client.

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { deviceAuthorizationClient } from "better-auth/client/plugins"; 
    
    export const authClient = createAuthClient({
      plugins: [
        deviceAuthorizationClient(), 
      ],
    });

The device flow follows these steps:

1.  **Device requests codes**: The device requests a device code and user code from the authorization server
2.  **User authorizes**: The user visits a verification URL and enters the user code
3.  **Device polls for token**: The device polls the server until the user completes authorization
4.  **Access granted**: Once authorized, the device receives an access token

To initiate device authorization, call `device.code` with the client ID:

    const { data, error } = await authClient.device.code({    client_id, // required    scope,});

| Prop | Description | Type |
| --- | --- | --- |
| `client_id` | 
The OAuth client identifier

 | `string;` |
| `scope?` | 

Space-separated list of requested scopes (optional)

 | `string;` |

Example usage:

    const { data } = await authClient.device.code({
      client_id: "your-client-id",
      scope: "openid profile email",
    });
    
    if (data) {
      console.log(`User code: ${data.user_code}`);
      console.log(`Verification URL: ${data.verification_uri}`);
      console.log(`Complete verification URL: ${data.verification_uri_complete}`);
    }

### [Polling for Token](#polling-for-token)

After displaying the user code, poll for the access token:

    const { data, error } = await authClient.device.token({    grant_type, // required    device_code, // required    client_id, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `grant_type` | 
Must be "urn:ietf:params:oauth:grant-type:device\_code"

 | `string;` |
| `device_code` | 

The device code from the initial request

 | `string;` |
| `client_id` | 

The OAuth client identifier

 | `string;` |

Example polling implementation:

    let pollingInterval = 5; // Start with 5 seconds
    const pollForToken = async () => {
      const { data, error } = await authClient.device.token({
        grant_type: "urn:ietf:params:oauth:grant-type:device_code",
        device_code,
        client_id: yourClientId,
        fetchOptions: {
          headers: {
            "user-agent": `My CLI`,
          },
        },
      });
    
      if (data?.access_token) {
        console.log("Authorization successful!");
      } else if (error) {
        switch (error.error) {
          case "authorization_pending":
            // Continue polling
            break;
          case "slow_down":
            pollingInterval += 5;
            break;
          case "access_denied":
            console.error("Access was denied by the user");
            return;
          case "expired_token":
            console.error("The device code has expired. Please try again.");
            return;
          default:
            console.error(`Error: ${error.error_description}`);
            return;
        }
        setTimeout(pollForToken, pollingInterval * 1000);
      }
    };
    
    pollForToken();

### [User Authorization Flow](#user-authorization-flow)

The user authorization flow requires two steps:

1.  **Code Verification**: Check if the entered user code is valid
2.  **Authorization**: User must be authenticated to approve/deny the device

Users must be authenticated before they can approve or deny device authorization requests. If not authenticated, redirect them to the login page with a return URL.

Create a page where users can enter their code:

app/device/page.tsx

    export default function DeviceAuthorizationPage() {
      const [userCode, setUserCode] = useState("");
      const [error, setError] = useState(null);
      
      const handleSubmit = async (e) => {
        e.preventDefault();
        
        try {
          // Format the code: remove dashes and convert to uppercase
          const formattedCode = userCode.trim().replace(/-/g, "").toUpperCase();
    
          // Check if the code is valid using GET /device endpoint
          const response = await authClient.device({
            query: { user_code: formattedCode },
          });
          
          if (response.data) {
            // Redirect to approval page
            window.location.href = `/device/approve?user_code=${formattedCode}`;
          }
        } catch (err) {
          setError("Invalid or expired code");
        }
      };
      
      return (
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={userCode}
            onChange={(e) => setUserCode(e.target.value)}
            placeholder="Enter device code (e.g., ABCD-1234)"
            maxLength={12}
          />
          <button type="submit">Continue</button>
          {error && <p>{error}</p>}
        </form>
      );
    }

### [Approving or Denying Device](#approving-or-denying-device)

Users must be authenticated to approve or deny device authorization requests:

#### [Approve Device](#approve-device)

    const { data, error } = await authClient.device.approve({    userCode, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userCode` | 
The user code to approve

 | `string;` |

#### [Deny Device](#deny-device)

    const { data, error } = await authClient.device.deny({    userCode, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userCode` | 
The user code to deny

 | `string;` |

#### [Example Approval Page](#example-approval-page)

app/device/approve/page.tsx

    export default function DeviceApprovalPage() {
      const { user } = useAuth(); // Must be authenticated
      const searchParams = useSearchParams();
      const userCode = searchParams.get("userCode");
      const [isProcessing, setIsProcessing] = useState(false);
      
      const handleApprove = async () => {
        setIsProcessing(true);
        try {
          await authClient.device.approve({
            userCode: userCode,
          });
          // Show success message
          alert("Device approved successfully!");
          window.location.href = "/";
        } catch (error) {
          alert("Failed to approve device");
        }
        setIsProcessing(false);
      };
      
      const handleDeny = async () => {
        setIsProcessing(true);
        try {
          await authClient.device.deny({
            userCode: userCode,
          });
          alert("Device denied");
          window.location.href = "/";
        } catch (error) {
          alert("Failed to deny device");
        }
        setIsProcessing(false);
      };
    
      if (!user) {
        // Redirect to login if not authenticated
        window.location.href = `/login?redirect=/device/approve?user_code=${userCode}`;
        return null;
      }
      
      return (
        <div>
          <h2>Device Authorization Request</h2>
          <p>A device is requesting access to your account.</p>
          <p>Code: {userCode}</p>
          
          <button onClick={handleApprove} disabled={isProcessing}>
            Approve
          </button>
          <button onClick={handleDeny} disabled={isProcessing}>
            Deny
          </button>
        </div>
      );
    }

### [Client Validation](#client-validation)

You can validate client IDs to ensure only authorized applications can use the device flow:

    deviceAuthorization({
      validateClient: async (clientId) => {
        // Check if client is authorized
        const client = await db.oauth_clients.findOne({ id: clientId });
        return client && client.allowDeviceFlow;
      },
      
      onDeviceAuthRequest: async (clientId, scope) => {
        // Log device authorization requests
        await logDeviceAuthRequest(clientId, scope);
      },
    })

### [Custom Code Generation](#custom-code-generation)

Customize how device and user codes are generated:

    deviceAuthorization({
      generateDeviceCode: async () => {
        // Custom device code generation
        return crypto.randomBytes(32).toString("hex");
      },
      
      generateUserCode: async () => {
        // Custom user code generation
        // Default uses: ABCDEFGHJKLMNPQRSTUVWXYZ23456789
        // (excludes 0, O, 1, I to avoid confusion)
        const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789";
        let code = "";
        for (let i = 0; i < 8; i++) {
          code += charset[Math.floor(Math.random() * charset.length)];
        }
        return code;
      },
    })

The device flow defines specific error codes:

| Error Code | Description |
| --- | --- |
| `authorization_pending` | User hasn't approved yet (continue polling) |
| `slow_down` | Polling too frequently (increase interval) |
| `expired_token` | Device code has expired |
| `access_denied` | User denied the authorization |
| `invalid_grant` | Invalid device code or client ID |

Here's a complete example for a CLI application based on the actual demo:

cli-auth.ts

    import { createAuthClient } from "better-auth/client";
    import { deviceAuthorizationClient } from "better-auth/client/plugins";
    import open from "open";
    
    const authClient = createAuthClient({
      baseURL: "http://localhost:3000",
      plugins: [deviceAuthorizationClient()],
    });
    
    async function authenticateCLI() {
      console.log("🔐 Better Auth Device Authorization Demo");
      console.log("⏳ Requesting device authorization...");
      
      try {
        // Request device code
        const { data, error } = await authClient.device.code({
          client_id: "demo-cli",
          scope: "openid profile email",
        });
        
        if (error || !data) {
          console.error("❌ Error:", error?.error_description);
          process.exit(1);
        }
        
        const {
          device_code,
          user_code,
          verification_uri,
          verification_uri_complete,
          interval = 5,
        } = data;
        
        console.log("\n📱 Device Authorization in Progress");
        console.log(`Please visit: ${verification_uri}`);
        console.log(`Enter code: ${user_code}\n`);
        
        // Open browser to verification page
        const urlToOpen = verification_uri_complete || verification_uri;
        console.log("🌐 Opening browser...");
        await open(urlToOpen);
        
        console.log(`⏳ Waiting for authorization... (polling every ${interval}s)`);
        
        // Poll for token
        await pollForToken(device_code, interval);
      } catch (err) {
        console.error("❌ Error:", err.message);
        process.exit(1);
      }
    }
    
    async function pollForToken(deviceCode: string, interval: number) {
      let pollingInterval = interval;
      
      return new Promise<void>((resolve) => {
        const poll = async () => {
          try {
            const { data, error } = await authClient.device.token({
              grant_type: "urn:ietf:params:oauth:grant-type:device_code",
              device_code: deviceCode,
              client_id: "demo-cli",
            });
            
            if (data?.access_token) {
              console.log("\nAuthorization Successful!");
              console.log("Access token received!");
              
              // Get user session
              const { data: session } = await authClient.getSession({
                fetchOptions: {
                  headers: {
                    Authorization: `Bearer ${data.access_token}`,
                  },
                },
              });
              
              console.log(`Hello, ${session?.user?.name || "User"}!`);
              resolve();
              process.exit(0);
            } else if (error) {
              switch (error.error) {
                case "authorization_pending":
                  // Continue polling silently
                  break;
                case "slow_down":
                  pollingInterval += 5;
                  console.log(`⚠️  Slowing down polling to ${pollingInterval}s`);
                  break;
                case "access_denied":
                  console.error("❌ Access was denied by the user");
                  process.exit(1);
                  break;
                case "expired_token":
                  console.error("❌ The device code has expired. Please try again.");
                  process.exit(1);
                  break;
                default:
                  console.error("❌ Error:", error.error_description);
                  process.exit(1);
              }
            }
          } catch (err) {
            console.error("❌ Network error:", err.message);
            process.exit(1);
          }
          
          // Schedule next poll
          setTimeout(poll, pollingInterval * 1000);
        };
        
        // Start polling
        setTimeout(poll, pollingInterval * 1000);
      });
    }
    
    // Run the authentication flow
    authenticateCLI().catch((err) => {
      console.error("❌ Fatal error:", err);
      process.exit(1);
    });

1.  **Rate Limiting**: The plugin enforces polling intervals to prevent abuse
2.  **Code Expiration**: Device and user codes expire after the configured time (default: 30 minutes)
3.  **Client Validation**: Always validate client IDs in production to prevent unauthorized access
4.  **HTTPS Only**: Always use HTTPS in production for device authorization
5.  **User Code Format**: User codes use a limited character set (excluding similar-looking characters like 0/O, 1/I) to reduce typing errors
6.  **Authentication Required**: Users must be authenticated before they can approve or deny device requests

### [Server](#server)

**verificationUri**: The URL of the verification page where users can enter their device code. Match this to the route of your verification page. Returned as `verification_uri` in the response. Can be an absolute URL (e.g., `https://example.com/device`) or relative path (e.g., `/device`). Default: `/device`.

**expiresIn**: The expiration time for device codes. Default: `"30m"` (30 minutes).

**interval**: The minimum polling interval. Default: `"5s"` (5 seconds).

**userCodeLength**: The length of the user code. Default: `8`.

**deviceCodeLength**: The length of the device code. Default: `40`.

**generateDeviceCode**: Custom function to generate device codes. Returns a string or `Promise<string>`.

**generateUserCode**: Custom function to generate user codes. Returns a string or `Promise<string>`.

**validateClient**: Function to validate client IDs. Takes a clientId and returns boolean or `Promise<boolean>`.

**onDeviceAuthRequest**: Hook called when device authorization is requested. Takes clientId and optional scope.

### [Client](#client)

No client-specific configuration options. The plugin adds the following methods:

*   **device()**: Verify user code validity
*   **device.code()**: Request device and user codes
*   **device.token()**: Poll for access token
*   **device.approve()**: Approve device (requires authentication)
*   **device.deny()**: Deny device (requires authentication)

The plugin requires a new table to store device authorization data.

Table Name: `deviceCode`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for the device authorization request |
| deviceCode | string | \- | The device verification code |
| userCode | string | \- | The user-friendly code for verification |
| userId | string |  | The ID of the user who approved/denied |
| clientId | string |  | The OAuth client identifier |
| scope | string |  | Requested OAuth scopes |
| status | string | \- | Current status: pending, approved, or denied |
| expiresAt | Date | \- | When the device code expires |
| lastPolledAt | Date |  | Last time the device polled for status |
| pollingInterval | number |  | Minimum seconds between polls |
| createdAt | Date | \- | When the request was created |
| updatedAt | Date | \- | When the request was last updated |</content>
</page>

<page>
  <title>FAQ | Better Auth</title>
  <url>https://www.better-auth.com/docs/reference/faq</url>
  <content>This page contains frequently asked questions, common issues, and other helpful information about Better Auth.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/reference/faq.mdx)

[

Previous Page

Telemetry

](https://www.better-auth.com/docs/reference/telemetry)

### On this page

No Headings</content>
</page>

<page>
  <title>Create your first plugin | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/your-first-plugin</url>
  <content>In this guide, we’ll walk you through the steps of creating your first Better Auth plugin.

This guide assumes you have [setup the basics](https://www.better-auth.com/docs/installation) of Better Auth and are ready to create your first plugin.

Before beginning, you must know what plugin you intend to create.

In this guide, we’ll create a **birthday plugin** to keep track of user birth dates.

Better Auth plugins operate as a pair: a [server plugin](https://www.better-auth.com/docs/concepts/plugins#create-a-server-plugin) and a [client plugin](https://www.better-auth.com/docs/concepts/plugins#create-a-client-plugin). The server plugin forms the foundation of your authentication system, while the client plugin provides convenient frontend APIs to interact with your server implementation.

You can read more about server/client plugins in our [documentation](https://www.better-auth.com/docs/concepts/plugins#creating-a-plugin).

### [Creating the server plugin](#creating-the-server-plugin)

Go ahead and find a suitable location to create your birthday plugin folder, with an `index.ts` file within.

In the `index.ts` file, we’ll export a function that represents our server plugin. This will be what we will later add to our plugin list in the `auth.ts` file.

index.ts

    import { createAuthClient } from "better-auth/client";
    import type { BetterAuthPlugin } from "better-auth";
    
    export const birthdayPlugin = () =>
      ({
        id: "birthdayPlugin",
      } satisfies BetterAuthPlugin);

Although this does nothing, you have technically just made yourself your first plugin, congratulations! 🎉

### [Defining a schema](#defining-a-schema)

In order to save each user’s birthday data, we must create a schema on top of the `user` model.

By creating a schema here, this also allows [Better Auth’s CLI](https://www.better-auth.com/docs/concepts/cli) to generate the schemas required to update your database.

index.ts

    //...
    export const birthdayPlugin = () =>
      ({
        id: "birthdayPlugin",
        schema: {
          user: {
            fields: {
              birthday: {
                type: "date", // string, number, boolean, date
                required: true, // if the field should be required on a new record. (default: false)
                unique: false, // if the field should be unique. (default: false)
                references: null // if the field is a reference to another table. (default: null)
              },
            },
          },
        },
      } satisfies BetterAuthPlugin);

For this example guide, we’ll set up authentication logic to check and ensure that the user who signs-up is older than 5. But the same concept could be applied for something like verifying users agreeing to the TOS or anything alike.

To do this, we’ll utilize [Hooks](https://www.better-auth.com/docs/concepts/plugins#hooks), which allows us to run code `before` or `after` an action is performed.

index.ts

    export const birthdayPlugin = () => ({
        //...
        // In our case, we want to write authorization logic,
        // meaning we want to intercept it `before` hand.
        hooks: {
          before: [
            {
              matcher: (context) => /* ... */,
              handler: createAuthMiddleware(async (ctx) => {
                //...
              }),
            },
          ],
        },
    } satisfies BetterAuthPlugin)

In our case we want to match any requests going to the signup path:

Before hook

    {
      matcher: (context) => context.path.startsWith("/sign-up/email"),
      //...
    }

And for our logic, we’ll write the following code to check the if user’s birthday makes them above 5 years old.

Imports

    import { APIError } from "better-auth/api";
    import { createAuthMiddleware } from "better-auth/plugins";

Before hook

    {
      //...
      handler: createAuthMiddleware(async (ctx) => {
        const { birthday } = ctx.body;
        if(!(birthday instanceof Date)) {
          throw new APIError("BAD_REQUEST", { message: "Birthday must be of type Date." });
        }
    
        const today = new Date();
        const fiveYearsAgo = new Date(today.setFullYear(today.getFullYear() - 5));
    
        if(birthday >= fiveYearsAgo) {
          throw new APIError("BAD_REQUEST", { message: "User must be above 5 years old." });
        }
    
        return { context: ctx };
      }),
    }

**Authorized!** 🔒

We’ve now successfully written code to ensure authorization for users above 5!

We’re close to the finish line! 🏁

Now that we have created our server plugin, the next step is to develop our client plugin. Since there isn’t much frontend APIs going on for this plugin, there isn’t much to do!

First, let’s create our `client.ts` file first:

Then, add the following code:

client.ts

    import { BetterAuthClientPlugin } from "better-auth";
    import type { birthdayPlugin } from "./index"; // make sure to import the server plugin as a type
    
    type BirthdayPlugin = typeof birthdayPlugin;
    
    export const birthdayClientPlugin = () => {
      return {
        id: "birthdayPlugin",
        $InferServerPlugin: {} as ReturnType<BirthdayPlugin>,
      } satisfies BetterAuthClientPlugin;
    };

What we’ve done is allow the client plugin to infer the types defined by our schema from the server plugin.

And that’s it! This is all it takes for the birthday client plugin. 🎂

Both the `client` and `server` plugins are now ready, the last step is to import them to both your `auth-client.ts` and your `server.ts` files respectively to initiate the plugin.

### [Server initiation](#server-initiation)

server.ts

    import { betterAuth } from "better-auth";
    import { birthdayPlugin } from "./birthday-plugin";
     
    export const auth = betterAuth({
        plugins: [
          birthdayPlugin(),
        ]
    });

### [Client initiation](#client-initiation)

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { birthdayClientPlugin } from "./birthday-plugin/client";
     
    const authClient = createAuthClient({
        plugins: [
          birthdayClientPlugin()
        ]
    });

### [Oh yeah, the schemas!](#oh-yeah-the-schemas)

Don’t forget to add your `birthday` field to your `user` table model!

Or, use the `generate` [CLI command](https://www.better-auth.com/docs/concepts/cli#generate):

    npx @better-auth/cli@latest generate

Congratulations! You’ve successfully created your first ever Better Auth plugin. We highly recommend you visit our [plugins documentation](https://www.better-auth.com/docs/concepts/plugins) to learn more information.

If you have a plugin you’d like to share with the community, feel free to let us know through our [Discord server](https://discord.gg/better-auth), or through a [pull-request](https://github.com/better-auth/better-auth/pulls) and we may add it to the [community-plugins](https://www.better-auth.com/docs/plugins/community-plugins) list!</content>
</page>

<page>
  <title>Community Plugins | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/community-plugins</url>
  <content>This page showcases a list of recommended community made plugins.

To create your own custom plugin, get started by reading our [plugins documentation](https://www.better-auth.com/docs/concepts/plugins). And if you want to share your plugin with the community, please open a pull request to add it to this list.

| 
Plugin

 | Description | 

Author

 |
| --- | --- | --- |
| [@dymo-api/better-auth](https://github.com/TPEOficial/dymo-api-better-auth) | Sign Up Protection and validation of disposable emails (the world's largest database with nearly 14 million entries). | [TPEOficial](https://github.com/TPEOficial) |
| [better-auth-harmony](https://github.com/gekorm/better-auth-harmony/) | Email & phone normalization and additional validation, blocking over 55,000 temporary email domains. | [GeKorm](https://github.com/GeKorm) |
| [validation-better-auth](https://github.com/Daanish2003/validation-better-auth) | Validate API request using any validation library (e.g., Zod, Yup) | [Daanish2003](https://github.com/Daanish2003) |
| [better-auth-localization](https://github.com/marcellosso/better-auth-localization) | Localize and customize better-auth messages with easy translation and message override support. | [marcellosso](https://github.com/marcellosso) |
| [better-auth-attio-plugin](https://github.com/tobimori/better-auth-attio-plugin) | Sync your products Better Auth users & workspaces with Attio | [tobimori](https://github.com/tobimori) |
| [better-auth-cloudflare](https://github.com/zpg6/better-auth-cloudflare) | Seamlessly integrate with Cloudflare Workers, D1, Hyperdrive, KV, R2, and geolocation services. Includes CLI for project generation, automated resource provisioning on Cloudflare, and database migrations. Supports Next.js, Hono, and more! | [zpg6](https://github.com/zpg6) |
| [expo-better-auth-passkey](https://github.com/kevcube/expo-better-auth-passkey) | Better-auth client plugin for using passkeys on mobile platforms in expo apps. Supports iOS, macOS, Android (and web!) by wrapping the existing better-auth passkey client plugin. | [kevcube](https://github.com/kevcube) |
| [better-auth-credentials-plugin](https://github.com/erickweil/better-auth-credentials-plugin) | LDAP authentication plugin for Better Auth. | [erickweil](https://github.com/erickweil) |
| [better-auth-opaque](https://github.com/TheUntraceable/better-auth-opaque) | Provides database-breach resistant authentication using the zero-knowledge OPAQUE protocol. | [TheUntraceable](https://github.com/TheUntraceable) |
| [better-auth-firebase-auth](https://github.com/yultyyev/better-auth-firebase-auth) | Firebase Authentication plugin for Better Auth with built-in email service, Google Sign-In, and password reset functionality. | [yultyyev](https://github.com/yultyyev) |</content>
</page>

<page>
  <title>Migrating from Auth0 to Better Auth | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/auth0-migration-guide</url>
  <content>    import { ManagementClient } from 'auth0';
    import { generateRandomString, symmetricEncrypt } from "better-auth/crypto";
    import { auth } from '@/lib/auth';
    
    const auth0Client = new ManagementClient({
        domain: process.env.AUTH0_DOMAIN!,
        clientId: process.env.AUTH0_CLIENT_ID!,
        clientSecret: process.env.AUTH0_SECRET!,
    });
    
    
    
    function safeDateConversion(timestamp?: string | number): Date {
        if (!timestamp) return new Date();
    
        const numericTimestamp = typeof timestamp === 'string' ? Date.parse(timestamp) : timestamp;
    
        const milliseconds = numericTimestamp < 1000000000000 ? numericTimestamp * 1000 : numericTimestamp;
    
        const date = new Date(milliseconds);
    
        if (isNaN(date.getTime())) {
            console.warn(`Invalid timestamp: ${timestamp}, falling back to current date`);
            return new Date();
        }
    
        // Check for unreasonable dates (before 2000 or after 2100)
        const year = date.getFullYear();
        if (year < 2000 || year > 2100) {
            console.warn(`Suspicious date year: ${year}, falling back to current date`);
            return new Date();
        }
    
        return date;
    }
    
    // Helper function to generate backup codes for 2FA
    async function generateBackupCodes(secret: string) {
        const key = secret;
        const backupCodes = Array.from({ length: 10 })
            .fill(null)
            .map(() => generateRandomString(10, "a-z", "0-9", "A-Z"))
            .map((code) => `${code.slice(0, 5)}-${code.slice(5)}`);
    
        const encCodes = await symmetricEncrypt({
            data: JSON.stringify(backupCodes),
            key: key,
        });
        return encCodes;
    }
    
    function mapAuth0RoleToBetterAuthRole(auth0Roles: string[]) {
        if (typeof auth0Roles === 'string') return auth0Roles;
        if (Array.isArray(auth0Roles)) return auth0Roles.join(',');
    }
    // helper function to migrate password from auth0 to better auth for custom hashes and algs
    async function migratePassword(auth0User: any) {
        if (auth0User.password_hash) {
            if (auth0User.password_hash.startsWith('$2a$') || auth0User.password_hash.startsWith('$2b$')) {
                return auth0User.password_hash;
            }
        }
    
        if (auth0User.custom_password_hash) {
            const customHash = auth0User.custom_password_hash;
    
            if (customHash.algorithm === 'bcrypt') {
                const hash = customHash.hash.value;
                if (hash.startsWith('$2a$') || hash.startsWith('$2b$')) {
                    return hash;
                }
            }
    
            return JSON.stringify({
                algorithm: customHash.algorithm,
                hash: {
                    value: customHash.hash.value,
                    encoding: customHash.hash.encoding || 'utf8',
                    ...(customHash.hash.digest && { digest: customHash.hash.digest }),
                    ...(customHash.hash.key && {
                        key: {
                            value: customHash.hash.key.value,
                            encoding: customHash.hash.key.encoding || 'utf8'
                        }
                    })
                },
                ...(customHash.salt && {
                    salt: {
                        value: customHash.salt.value,
                        encoding: customHash.salt.encoding || 'utf8',
                        position: customHash.salt.position || 'prefix'
                    }
                }),
                ...(customHash.password && {
                    password: {
                        encoding: customHash.password.encoding || 'utf8'
                    }
                }),
                ...(customHash.algorithm === 'scrypt' && {
                    keylen: customHash.keylen,
                    cost: customHash.cost || 16384,
                    blockSize: customHash.blockSize || 8,
                    parallelization: customHash.parallelization || 1
                })
            });
        }
    
        return null;
    }
    
    async function migrateMFAFactors(auth0User: any, userId: string | undefined, ctx: any) {
        if (!userId || !auth0User.mfa_factors || !Array.isArray(auth0User.mfa_factors)) {
            return;
        }
    
        for (const factor of auth0User.mfa_factors) {
            try {
                if (factor.totp && factor.totp.secret) {
                    await ctx.adapter.create({
                        model: "twoFactor",
                        data: {
                            userId: userId,
                            secret: factor.totp.secret,
                            backupCodes: await generateBackupCodes(factor.totp.secret)
                        }
                    });
                }
            } catch (error) {
                console.error(`Failed to migrate MFA factor for user ${userId}:`, error);
            }
        }
    }
    
    async function migrateOAuthAccounts(auth0User: any, userId: string | undefined, ctx: any) {
        if (!userId || !auth0User.identities || !Array.isArray(auth0User.identities)) {
            return;
        }
    
        for (const identity of auth0User.identities) {
            try {
                const providerId = identity.provider === 'auth0' ? "credential" : identity.provider.split("-")[0];
                await ctx.adapter.create({
                    model: "account",
                    data: {
                        id: `${auth0User.user_id}|${identity.provider}|${identity.user_id}`,
                        userId: userId,
                        password: await migratePassword(auth0User),
                        providerId: providerId || identity.provider,
                        accountId: identity.user_id,
                        accessToken: identity.access_token,
                        tokenType: identity.token_type,
                        refreshToken: identity.refresh_token,
                        accessTokenExpiresAt: identity.expires_in ? new Date(Date.now() + identity.expires_in * 1000) : undefined,
                        // if you are enterprise user, you can get the refresh tokens or all the tokensets - auth0Client.users.getAllTokensets 
                        refreshTokenExpiresAt: identity.refresh_token_expires_in ? new Date(Date.now() + identity.refresh_token_expires_in * 1000) : undefined,
    
                        scope: identity.scope,
                        idToken: identity.id_token,
                        createdAt: safeDateConversion(auth0User.created_at),
                        updatedAt: safeDateConversion(auth0User.updated_at)
                    },
                    forceAllowId: true
                }).catch((error: Error) => {
                    console.error(`Failed to create OAuth account for user ${userId} with provider ${providerId}:`, error);
                    return ctx.adapter.create({
                        // Try creating without optional fields if the first attempt failed
                        model: "account",
                        data: {
                            id: `${auth0User.user_id}|${identity.provider}|${identity.user_id}`,
                            userId: userId,
                            password: migratePassword(auth0User),
                            providerId: providerId,
                            accountId: identity.user_id,
                            accessToken: identity.access_token,
                            tokenType: identity.token_type,
                            refreshToken: identity.refresh_token,
                            accessTokenExpiresAt: identity.expires_in ? new Date(Date.now() + identity.expires_in * 1000) : undefined,
                            refreshTokenExpiresAt: identity.refresh_token_expires_in ? new Date(Date.now() + identity.refresh_token_expires_in * 1000) : undefined,
                            scope: identity.scope,
                            idToken: identity.id_token,
                            createdAt: safeDateConversion(auth0User.created_at),
                            updatedAt: safeDateConversion(auth0User.updated_at)
                        },
                        forceAllowId: true
                    });
                });
    
                console.log(`Successfully migrated OAuth account for user ${userId} with provider ${providerId}`);
            } catch (error) {
                console.error(`Failed to migrate OAuth account for user ${userId}:`, error);
            }
        }
    }
    
    async function migrateOrganizations(ctx: any) {
        try {
            const organizations = await auth0Client.organizations.getAll();
            for (const org of organizations.data || []) {
                try {
                    await ctx.adapter.create({
                        model: "organization",
                        data: {
                            id: org.id,
                            name: org.display_name || org.id,
                            slug: (org.display_name || org.id).toLowerCase().replace(/[^a-z0-9]/g, '-'),
                            logo: org.branding?.logo_url,
                            metadata: JSON.stringify(org.metadata || {}),
                            createdAt: safeDateConversion(org.created_at),
                        },
                        forceAllowId: true
                    });
                    const members = await auth0Client.organizations.getMembers({ id: org.id });
                    for (const member of members.data || []) {
                        try {
                            const userRoles = await auth0Client.organizations.getMemberRoles({
                                id: org.id,
                                user_id: member.user_id
                            });
                            const role = mapAuth0RoleToBetterAuthRole(userRoles.data?.map(r => r.name) || []);
                            await ctx.adapter.create({
                                model: "member",
                                data: {
                                    id: `${org.id}|${member.user_id}`,
                                    organizationId: org.id,
                                    userId: member.user_id,
                                    role: role,
                                    createdAt: new Date()
                                },
                                forceAllowId: true
                            });
    
                            console.log(`Successfully migrated member ${member.user_id} for organization ${org.display_name || org.id}`);
                        } catch (error) {
                            console.error(`Failed to migrate member ${member.user_id} for organization ${org.display_name || org.id}:`, error);
                        }
                    }
    
                    console.log(`Successfully migrated organization: ${org.display_name || org.id}`);
                } catch (error) {
                    console.error(`Failed to migrate organization ${org.display_name || org.id}:`, error);
                }
            }
            console.log('Organization migration completed');
        } catch (error) {
            console.error('Failed to migrate organizations:', error);
        }
    }
    
    async function migrateFromAuth0() {
        try {
            const ctx = await auth.$context;
            const isAdminEnabled = ctx.options?.plugins?.find(plugin => plugin.id === "admin");
            const isUsernameEnabled = ctx.options?.plugins?.find(plugin => plugin.id === "username");
            const isOrganizationEnabled = ctx.options?.plugins?.find(plugin => plugin.id === "organization");
            const perPage = 100;
            const auth0Users: any[] = [];
            let pageNumber = 0;
    
            while (true) {
                try {
                    const params = {
                        per_page: perPage,
                        page: pageNumber,
                        include_totals: true,
                    };
                    const response = (await auth0Client.users.getAll(params)).data as any;
                    const users = response.users || [];
                    if (users.length === 0) break;
                    auth0Users.push(...users);
                    pageNumber++;
    
                    if (users.length < perPage) break;
                } catch (error) {
                    console.error('Error fetching users:', error);
                    break;
                }
            }
    
    
            console.log(`Found ${auth0Users.length} users to migrate`);
    
            for (const auth0User of auth0Users) {
                try {
                    // Determine if this is a password-based or OAuth user
                    const isOAuthUser = auth0User.identities?.some((identity: any) => identity.provider !== 'auth0');
                    // Base user data that's common for both types
                    const baseUserData = {
                        id: auth0User.user_id,
                        email: auth0User.email,
                        emailVerified: auth0User.email_verified || false,
                        name: auth0User.name || auth0User.nickname,
                        image: auth0User.picture,
                        createdAt: safeDateConversion(auth0User.created_at),
                        updatedAt: safeDateConversion(auth0User.updated_at),
                        ...(isAdminEnabled ? {
                            banned: auth0User.blocked || false,
                            role: mapAuth0RoleToBetterAuthRole(auth0User.roles || []),
                        } : {}),
    
                        ...(isUsernameEnabled ? {
                            username: auth0User.username || auth0User.nickname,
                        } : {}),
    
                    };
    
                    const createdUser = await ctx.adapter.create({
                        model: "user",
                        data: {
                            ...baseUserData,
                        },
                        forceAllowId: true
                    });
    
                    if (!createdUser?.id) {
                        throw new Error('Failed to create user');
                    }
    
    
                    await migrateOAuthAccounts(auth0User, createdUser.id, ctx)
                    console.log(`Successfully migrated user: ${auth0User.email}`);
                } catch (error) {
                    console.error(`Failed to migrate user ${auth0User.email}:`, error);
                }
            }
            if (isOrganizationEnabled) {
                await migrateOrganizations(ctx);
            }
            // the reset of migration will be here.
            console.log('Migration completed successfully');
        } catch (error) {
            console.error('Migration failed:', error);
            throw error;
        }
    }
    
    migrateFromAuth0()
        .then(() => {
            console.log('Migration completed');
            process.exit(0);
        })
        .catch((error) => {
            console.error('Migration failed:', error);
            process.exit(1);
        });</content>
</page>

<page>
  <title>Dub | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/dub</url>
  <content>[Dub](https://dub.co/) is an open source modern link management platform for entrepreneurs, creators, and growth teams.

This plugins allows you to track leads when a user signs up using a Dub link. It also adds OAuth linking support to allow you to build integrations extending Dub's linking management infrastructure.

### [Install the plugin](#install-the-plugin)

First, install the plugin:

### [Install the Dub SDK](#install-the-dub-sdk)

Next, install the Dub SDK on your server:

### [Configure the plugin](#configure-the-plugin)

Add the plugin to your auth config:

auth.ts

    import { betterAuth } from "better-auth"
    import { dubAnalytics } from "@dub/better-auth"
    import { dub } from "dub"
    
    export const auth = betterAuth({
        plugins: [
            dubAnalytics({
                dubClient: new Dub()
            })
        ]
    })

### [Lead Tracking](#lead-tracking)

By default, the plugin will track sign up events as leads. You can disable this by setting `disableLeadTracking` to `true`.

    import { dubAnalytics } from "@dub/better-auth";
    import { betterAuth } from "better-auth";
    import { Dub } from "dub";
    
    const dub = new Dub();
    
    const betterAuth = betterAuth({
      plugins: [
        dubAnalytics({
          dubClient: dub,
          disableLeadTracking: true, // Disable lead tracking
        }),
      ],
    });

### [OAuth Linking](#oauth-linking)

The plugin supports OAuth for account linking.

First, you need to setup OAuth app in Dub. Dub supports OAuth 2.0 authentication, which is recommended if you build integrations extending Dub’s functionality [Learn more about OAuth](https://dub.co/docs/integrations/quickstart#integrating-via-oauth-2-0-recommended).

Once you get the client ID and client secret, you can configure the plugin.

    dubAnalytics({
      dubClient: dub,
      oauth: {
        clientId: "your-client-id",
        clientSecret: "your-client-secret",
      },
    });

And in the client, you need to use the `dubAnalyticsClient` plugin.

    import { createAuthClient } from "better-auth/client";
    import { dubAnalyticsClient } from "@dub/better-auth/client";
    
    const authClient = createAuthClient({
      plugins: [dubAnalyticsClient()],
    });

To link account with Dub, you need to use the `dub.link`.

    const { data, error } = await authClient.dub.link({    callbackURL: "/dashboard", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `callbackURL` | 
URL to redirect to after linking

 | `string` |

You can pass the following options to the plugin:

### [`dubClient`](#dubclient)

The Dub client instance.

### [`disableLeadTracking`](#disableleadtracking)

Disable lead tracking for sign up events.

### [`leadEventName`](#leadeventname)

Event name for sign up leads.

### [`customLeadTrack`](#customleadtrack)

Custom lead track function.

### [`oauth`](#oauth)

Dub OAuth configuration.

### [`oauth.clientId`](#oauthclientid)

Client ID for Dub OAuth.

### [`oauth.clientSecret`](#oauthclientsecret)

Client secret for Dub OAuth.

### [`oauth.pkce`](#oauthpkce)

Enable PKCE for Dub OAuth.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/dub.mdx)</content>
</page>

<page>
  <title>Admin | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/admin</url>
  <content>The Admin plugin provides a set of administrative functions for user management in your application. It allows administrators to perform various operations such as creating users, managing user roles, banning/unbanning users, impersonating users, and more.

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

To use the Admin plugin, add it to your auth config.

auth.ts

    import { betterAuth } from "better-auth"
    import { admin } from "better-auth/plugins"
    
    export const auth = betterAuth({
        // ... other config options
        plugins: [
            admin() 
        ]
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

Next, include the admin client plugin in your authentication client instance.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { adminClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            adminClient()
        ]
    })

Before performing any admin operations, the user must be authenticated with an admin account. An admin is any user assigned the `admin` role or any user whose ID is included in the `adminUserIds` option.

### [Create User](#create-user)

Allows an admin to create a new user.

    const { data: newUser, error } = await authClient.admin.createUser({    email: "user@example.com", // required    password: "some-secure-password", // required    name: "James Smith", // required    role: "user",    data: { customField: "customValue" },});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
The email of the user.

 | `string` |
| `password` | 

The password of the user.

 | `string` |
| `name` | 

The name of the user.

 | `string` |
| `role?` | 

A string or array of strings representing the roles to apply to the new user.

 | `string | string[]` |
| `data?` | 

Extra fields for the user. Including custom additional fields.

 | `Record<string, any>` |

### [List Users](#list-users)

Allows an admin to list all users in the database.

Notes

All properties are optional to configure. By default, 100 rows are returned, you can configure this by the `limit` property.

    const { data: users, error } = await authClient.admin.listUsers({    query: {        searchValue: "some name",        searchField: "name",        searchOperator: "contains",        limit: 100,        offset: 100,        sortBy: "name",        sortDirection: "desc",        filterField: "email",        filterValue: "hello@example.com",        filterOperator: "eq",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `searchValue?` | 
The value to search for.

 | `string` |
| `searchField?` | 

The field to search in, defaults to email. Can be `email` or `name`.

 | `"email" | "name"` |
| `searchOperator?` | 

The operator to use for the search. Can be `contains`, `starts_with` or `ends_with`.

 | `"contains" | "starts_with" | "ends_with"` |
| `limit?` | 

The number of users to return. Defaults to 100.

 | `string | number` |
| `offset?` | 

The offset to start from.

 | `string | number` |
| `sortBy?` | 

The field to sort by.

 | `string` |
| `sortDirection?` | 

The direction to sort by.

 | `"asc" | "desc"` |
| `filterField?` | 

The field to filter by.

 | `string` |
| `filterValue?` | 

The value to filter by.

 | `string | number | boolean` |
| `filterOperator?` | 

The operator to use for the filter.

 | `"eq" | "ne" | "lt" | "lte" | "gt" | "gte"` |

#### [Query Filtering](#query-filtering)

The `listUsers` function supports various filter operators including `eq`, `contains`, `starts_with`, and `ends_with`.

#### [Pagination](#pagination)

The `listUsers` function supports pagination by returning metadata alongside the user list. The response includes the following fields:

    {
      users: User[],   // Array of returned users
      total: number,   // Total number of users after filters and search queries
      limit: number | undefined,   // The limit provided in the query
      offset: number | undefined   // The offset provided in the query
    }

##### [How to Implement Pagination](#how-to-implement-pagination)

To paginate results, use the `total`, `limit`, and `offset` values to calculate:

*   **Total pages:** `Math.ceil(total / limit)`
*   **Current page:** `(offset / limit) + 1`
*   **Next page offset:** `Math.min(offset + limit, (total - 1))` – The value to use as `offset` for the next page, ensuring it does not exceed the total number of pages.
*   **Previous page offset:** `Math.max(0, offset - limit)` – The value to use as `offset` for the previous page (ensuring it doesn’t go below zero).

##### [Example Usage](#example-usage)

Fetching the second page with 10 users per page:

admin.ts

    const pageSize = 10;
    const currentPage = 2;
    
    const users = await authClient.admin.listUsers({
        query: {
            limit: pageSize,
            offset: (currentPage - 1) * pageSize
        }
    });
    
    const totalUsers = users.total;
    const totalPages = Math.ceil(totalUsers / pageSize)

### [Set User Role](#set-user-role)

Changes the role of a user.

    const { data, error } = await authClient.admin.setRole({    userId: "user-id",    role: "admin", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userId?` | 
The user id which you want to set the role for.

 | `string` |
| `role` | 

The role to set, this can be a string or an array of strings.

 | `string | string[]` |

### [Set User Password](#set-user-password)

Changes the password of a user.

POST

/admin/set-user-password

    const { data, error } = await authClient.admin.setUserPassword({    newPassword: 'new-password', // required    userId: 'user-id', // required});

| Prop | Description | Type |
| --- | --- | --- |
| `newPassword` | 
The new password.

 | `string` |
| `userId` | 

The user id which you want to set the password for.

 | `string` |

### [Update user](#update-user)

Update a user's details.

    const { data, error } = await authClient.admin.updateUser({    userId: "user-id", // required    data: { name: "John Doe" }, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
The user id which you want to update.

 | `string` |
| `data` | 

The data to update.

 | `Record<string, any>` |

### [Ban User](#ban-user)

Bans a user, preventing them from signing in and revokes all of their existing sessions.

    await authClient.admin.banUser({    userId: "user-id", // required    banReason: "Spamming",    banExpiresIn: 60 * 60 * 24 * 7,});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
The user id which you want to ban.

 | `string` |
| `banReason?` | 

The reason for the ban.

 | `string` |
| `banExpiresIn?` | 

The number of seconds until the ban expires. If not provided, the ban will never expire.

 | `number` |

### [Unban User](#unban-user)

Removes the ban from a user, allowing them to sign in again.

    await authClient.admin.unbanUser({    userId: "user-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
The user id which you want to unban.

 | `string` |

### [List User Sessions](#list-user-sessions)

Lists all sessions for a user.

POST

/admin/list-user-sessions

    const { data, error } = await authClient.admin.listUserSessions({    userId: "user-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
The user id.

 | `string` |

### [Revoke User Session](#revoke-user-session)

Revokes a specific session for a user.

POST

/admin/revoke-user-session

    const { data, error } = await authClient.admin.revokeUserSession({    sessionToken: "session_token_here", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `sessionToken` | 
The session token which you want to revoke.

 | `string` |

### [Revoke All Sessions for a User](#revoke-all-sessions-for-a-user)

Revokes all sessions for a user.

POST

/admin/revoke-user-sessions

    const { data, error } = await authClient.admin.revokeUserSessions({    userId: "user-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
The user id which you want to revoke all sessions for.

 | `string` |

### [Impersonate User](#impersonate-user)

This feature allows an admin to create a session that mimics the specified user. The session will remain active until either the browser session ends or it reaches 1 hour. You can change this duration by setting the `impersonationSessionDuration` option.

POST

/admin/impersonate-user

    const { data, error } = await authClient.admin.impersonateUser({    userId: "user-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
The user id which you want to impersonate.

 | `string` |

### [Stop Impersonating User](#stop-impersonating-user)

To stop impersonating a user and continue with the admin account, you can use `stopImpersonating`

POST

/admin/stop-impersonating

    await authClient.admin.stopImpersonating();

### [Remove User](#remove-user)

Hard deletes a user from the database.

    const { data: deletedUser, error } = await authClient.admin.removeUser({    userId: "user-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `userId` | 
The user id which you want to remove.

 | `string` |

The admin plugin offers a highly flexible access control system, allowing you to manage user permissions based on their role. You can define custom permission sets to fit your needs.

### [Roles](#roles)

By default, there are two roles:

`admin`: Users with the admin role have full control over other users.

`user`: Users with the user role have no control over other users.

A user can have multiple roles. Multiple roles are stored as string separated by comma (",").

### [Permissions](#permissions)

By default, there are two resources with up to six permissions.

**user**: `create` `list` `set-role` `ban` `impersonate` `delete` `set-password`

**session**: `list` `revoke` `delete`

Users with the admin role have full control over all the resources and actions. Users with the user role have no control over any of those actions.

### [Custom Permissions](#custom-permissions)

The plugin provides an easy way to define your own set of permissions for each role.

#### [Create Access Control](#create-access-control)

You first need to create an access controller by calling the `createAccessControl` function and passing the statement object. The statement object should have the resource name as the key and the array of actions as the value.

permissions.ts

    import { createAccessControl } from "better-auth/plugins/access";
    
    /**
     * make sure to use `as const` so typescript can infer the type correctly
     */
    const statement = { 
        project: ["create", "share", "update", "delete"], 
    } as const; 
    
    const ac = createAccessControl(statement); 

#### [Create Roles](#create-roles)

Once you have created the access controller you can create roles with the permissions you have defined.

permissions.ts

    import { createAccessControl } from "better-auth/plugins/access";
    
    export const statement = {
        project: ["create", "share", "update", "delete"], // <-- Permissions available for created roles
    } as const;
    
    const ac = createAccessControl(statement);
    
    export const user = ac.newRole({ 
        project: ["create"], 
    }); 
    
    export const admin = ac.newRole({ 
        project: ["create", "update"], 
    }); 
    
    export const myCustomRole = ac.newRole({ 
        project: ["create", "update", "delete"], 
        user: ["ban"], 
    }); 

When you create custom roles for existing roles, the predefined permissions for those roles will be overridden. To add the existing permissions to the custom role, you need to import `defaultStatements` and merge it with your new statement, plus merge the roles' permissions set with the default roles.

permissions.ts

    import { createAccessControl } from "better-auth/plugins/access";
    import { defaultStatements, adminAc } from "better-auth/plugins/admin/access";
    
    const statement = {
        ...defaultStatements, 
        project: ["create", "share", "update", "delete"],
    } as const;
    
    const ac = createAccessControl(statement);
    
    const admin = ac.newRole({
        project: ["create", "update"],
        ...adminAc.statements, 
    });

#### [Pass Roles to the Plugin](#pass-roles-to-the-plugin)

Once you have created the roles you can pass them to the admin plugin both on the client and the server.

auth.ts

    import { betterAuth } from "better-auth"
    import { admin as adminPlugin } from "better-auth/plugins"
    import { ac, admin, user } from "@/auth/permissions"
    
    export const auth = betterAuth({
        plugins: [
            adminPlugin({
                ac,
                roles: {
                    admin,
                    user,
                    myCustomRole
                }
            }),
        ],
    });

You also need to pass the access controller and the roles to the client plugin.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { adminClient } from "better-auth/client/plugins"
    import { ac, admin, user, myCustomRole } from "@/auth/permissions"
    
    export const client = createAuthClient({
        plugins: [
            adminClient({
                ac,
                roles: {
                    admin,
                    user,
                    myCustomRole
                }
            })
        ]
    })

### [Access Control Usage](#access-control-usage)

**Has Permission**:

To check a user's permissions, you can use the `hasPermission` function provided by the client.

POST

/admin/has-permission

    const { data, error } = await authClient.admin.hasPermission({    userId: "user-id",    permission: { "project": ["create", "update"] } /* Must use this, or permissions */,    permissions,});

| Prop | Description | Type |
| --- | --- | --- |
| `userId?` | 
The user id which you want to check the permissions for.

 | `string` |
| `permission?` | 

Optionally check if a single permission is granted. Must use this, or permissions.

 | `Record<string, string[]>` |
| `permissions?` | 

Optionally check if multiple permissions are granted. Must use this, or permission.

 | `Record<string, string[]>` |

Example usage:

auth-client.ts

    const canCreateProject = await authClient.admin.hasPermission({
      permissions: {
        project: ["create"],
      },
    });
    
    // You can also check multiple resource permissions at the same time
    const canCreateProjectAndCreateSale = await authClient.admin.hasPermission({
      permissions: {
        project: ["create"],
        sale: ["create"]
      },
    });

If you want to check a user's permissions server-side, you can use the `userHasPermission` action provided by the `api` to check the user's permissions.

api.ts

    import { auth } from "@/auth";
    
    await auth.api.userHasPermission({
      body: {
        userId: 'id', //the user id
        permissions: {
          project: ["create"], // This must match the structure in your access control
        },
      },
    });
    
    // You can also just pass the role directly
    await auth.api.userHasPermission({
      body: {
       role: "admin",
        permissions: {
          project: ["create"], // This must match the structure in your access control
        },
      },
    });
    
    // You can also check multiple resource permissions at the same time
    await auth.api.userHasPermission({
      body: {
       role: "admin",
        permissions: {
          project: ["create"], // This must match the structure in your access control
          sale: ["create"]
        },
      },
    });

**Check Role Permission**:

Use the `checkRolePermission` function on the client side to verify whether a given **role** has a specific **permission**. This is helpful after defining roles and their permissions, as it allows you to perform permission checks without needing to contact the server.

Note that this function does **not** check the permissions of the currently logged-in user directly. Instead, it checks what permissions are assigned to a specified role. The function is synchronous, so you don't need to use `await` when calling it.

auth-client.ts

    const canCreateProject = authClient.admin.checkRolePermission({
      permissions: {
        user: ["delete"],
      },
      role: "admin",
    });
    
    // You can also check multiple resource permissions at the same time
    const canDeleteUserAndRevokeSession = authClient.admin.checkRolePermission({
      permissions: {
        user: ["delete"],
        session: ["revoke"]
      },
      role: "admin",
    });

This plugin adds the following fields to the `user` table:

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| role | string |  | The user's role. Defaults to \`user\`. Admins will have the \`admin\` role. |
| banned | boolean |  | Indicates whether the user is banned. |
| banReason | string |  | The reason for the user's ban. |
| banExpires | date |  | The date when the user's ban will expire. |

And adds one field in the `session` table:

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| impersonatedBy | string |  | The ID of the admin that is impersonating this session. |

### [Default Role](#default-role)

The default role for a user. Defaults to `user`.

auth.ts

    admin({
      defaultRole: "regular",
    });

### [Admin Roles](#admin-roles)

The roles that are considered admin roles when **not** using custom access control. Defaults to `["admin"]`.

auth.ts

    admin({
      adminRoles: ["admin", "superadmin"],
    });

**Note:** The `adminRoles` option is **not required** when using custom access control (via `ac` and `roles`). When you define custom roles with specific permissions, those roles will have exactly the permissions you grant them through the access control system.

**Warning:** When **not** using custom access control, any role that isn't in the `adminRoles` list will **not** be able to perform admin operations.

### [Admin userIds](#admin-userids)

You can pass an array of userIds that should be considered as admin. Default to `[]`

auth.ts

    admin({
        adminUserIds: ["user_id_1", "user_id_2"]
    })

If a user is in the `adminUserIds` list, they will be able to perform any admin operation.

### [impersonationSessionDuration](#impersonationsessionduration)

The duration of the impersonation session in seconds. Defaults to 1 hour.

auth.ts

    admin({
      impersonationSessionDuration: 60 * 60 * 24, // 1 day
    });

### [Default Ban Reason](#default-ban-reason)

The default ban reason for a user created by the admin. Defaults to `No reason`.

auth.ts

    admin({
      defaultBanReason: "Spamming",
    });

### [Default Ban Expires In](#default-ban-expires-in)

The default ban expires in for a user created by the admin in seconds. Defaults to `undefined` (meaning the ban never expires).

auth.ts

    admin({
      defaultBanExpiresIn: 60 * 60 * 24, // 1 day
    });

### [bannedUserMessage](#bannedusermessage)

The message to show when a banned user tries to sign in. Defaults to "You have been banned from this application. Please contact support if you believe this is an error."

auth.ts

    admin({
      bannedUserMessage: "Custom banned user message",
    });</content>
</page>

<page>
  <title>Migrating from Clerk to Better Auth | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/clerk-migration-guide</url>
  <content>    import { generateRandomString, symmetricEncrypt } from "better-auth/crypto";
    
    import { auth } from "@/lib/auth"; // import your auth instance
    
    function getCSVData(csv: string) {
      const lines = csv.split('\n').filter(line => line.trim());
      const headers = lines[0]?.split(',').map(header => header.trim()) || [];
      const jsonData = lines.slice(1).map(line => {
          const values = line.split(',').map(value => value.trim());
          return headers.reduce((obj, header, index) => {
              obj[header] = values[index] || '';
              return obj;
          }, {} as Record<string, string>);
      });
    
      return jsonData as Array<{
          id: string;
          first_name: string;
          last_name: string;
          username: string;
          primary_email_address: string;
          primary_phone_number: string;
          verified_email_addresses: string;
          unverified_email_addresses: string;
          verified_phone_numbers: string;
          unverified_phone_numbers: string;
          totp_secret: string;
          password_digest: string;
          password_hasher: string;
      }>;
    }
    
    const exportedUserCSV = await Bun.file("exported_users.csv").text(); // this is the file you downloaded from Clerk
    
    async function getClerkUsers(totalUsers: number) {
      const clerkUsers: {
          id: string;
          first_name: string;
          last_name: string;
          username: string;
          image_url: string;
          password_enabled: boolean;
          two_factor_enabled: boolean;
          totp_enabled: boolean;
          backup_code_enabled: boolean;
          banned: boolean;
          locked: boolean;
          lockout_expires_in_seconds: number;
          created_at: number;
          updated_at: number;
          external_accounts: {
              id: string;
              provider: string;
              identification_id: string;
              provider_user_id: string;
              approved_scopes: string;
              email_address: string;
              first_name: string;
              last_name: string;
              image_url: string;
              created_at: number;
              updated_at: number;
          }[]
      }[] = [];
      for (let i = 0; i < totalUsers; i += 500) {
          const response = await fetch(`https://api.clerk.com/v1/users?offset=${i}&limit=${500}`, {
              headers: {
                  'Authorization': `Bearer ${process.env.CLERK_SECRET_KEY}`
              }
          });
          if (!response.ok) {
              throw new Error(`Failed to fetch users: ${response.statusText}`);
          }
          const clerkUsersData = await response.json();
          // biome-ignore lint/suspicious/noExplicitAny: <explanation>
          clerkUsers.push(...clerkUsersData as any);
      }
      return clerkUsers;
    }
    
    
    export async function generateBackupCodes(
      secret: string,
    ) {
      const key = secret;
      const backupCodes = Array.from({ length: 10 })
          .fill(null)
          .map(() => generateRandomString(10, "a-z", "0-9", "A-Z"))
          .map((code) => `${code.slice(0, 5)}-${code.slice(5)}`);
      const encCodes = await symmetricEncrypt({
          data: JSON.stringify(backupCodes),
          key: key,
      });
      return encCodes
    }
    
    // Helper function to safely convert timestamp to Date
    function safeDateConversion(timestamp?: number): Date {
      if (!timestamp) return new Date();
    
      const date = new Date(timestamp);
    
      // Check if the date is valid
      if (isNaN(date.getTime())) {
          console.warn(`Invalid timestamp: ${timestamp}, falling back to current date`);
          return new Date();
      }
    
      // Check for unreasonable dates (before 2000 or after 2100)
      const year = date.getFullYear();
      if (year < 2000 || year > 2100) {
          console.warn(`Suspicious date year: ${year}, falling back to current date`);
          return new Date();
      }
    
      return date;
    }
    
    async function migrateFromClerk() {
      const jsonData = getCSVData(exportedUserCSV);
      const clerkUsers = await getClerkUsers(jsonData.length);
      const ctx = await auth.$context
      const isAdminEnabled = ctx.options?.plugins?.find(plugin => plugin.id === "admin");
      const isTwoFactorEnabled = ctx.options?.plugins?.find(plugin => plugin.id === "two-factor");
      const isUsernameEnabled = ctx.options?.plugins?.find(plugin => plugin.id === "username");
      const isPhoneNumberEnabled = ctx.options?.plugins?.find(plugin => plugin.id === "phone-number");
      for (const user of jsonData) {
          const { id, first_name, last_name, username, primary_email_address, primary_phone_number, verified_email_addresses, unverified_email_addresses, verified_phone_numbers, unverified_phone_numbers, totp_secret, password_digest, password_hasher } = user;
          const clerkUser = clerkUsers.find(clerkUser => clerkUser?.id === id);
    
          // create user
          const createdUser = await ctx.adapter.create<{
              id: string;
          }>({
              model: "user",
              data: {
                  id,
                  email: primary_email_address,
                  emailVerified: verified_email_addresses.length > 0,
                  name: `${first_name} ${last_name}`,
                  image: clerkUser?.image_url,
                  createdAt: safeDateConversion(clerkUser?.created_at),
                  updatedAt: safeDateConversion(clerkUser?.updated_at),
                  // # Two Factor (if you enabled two factor plugin)
                  ...(isTwoFactorEnabled ? {
                      twoFactorEnabled: clerkUser?.two_factor_enabled
                  } : {}),
                  // # Admin (if you enabled admin plugin)
                  ...(isAdminEnabled ? {
                      banned: clerkUser?.banned,
                      banExpires: clerkUser?.lockout_expires_in_seconds
                         ? new Date(Date.now() + clerkUser.lockout_expires_in_seconds * 1000)
                         : undefined,
                      role: "user"
                  } : {}),
                  // # Username (if you enabled username plugin)
                  ...(isUsernameEnabled ? {
                      username: username,
                  } : {}),
                  // # Phone Number (if you enabled phone number plugin)  
                  ...(isPhoneNumberEnabled ? {
                      phoneNumber: primary_phone_number,
                      phoneNumberVerified: verified_phone_numbers.length > 0,
                  } : {}),
              },
              forceAllowId: true
          }).catch(async e => {
              return await ctx.adapter.findOne<{
                  id: string;
              }>({
                  model: "user",
                  where: [{
                      field: "id",
                      value: id
                  }]
              })
          })
          // create external account
          const externalAccounts = clerkUser?.external_accounts;
          if (externalAccounts) {
              for (const externalAccount of externalAccounts) {
                  const { id, provider, identification_id, provider_user_id, approved_scopes, email_address, first_name, last_name, image_url, created_at, updated_at } = externalAccount;
                  if (externalAccount.provider === "credential") {
                      await ctx.adapter.create({
                          model: "account",
                          data: {
                              id,
                              providerId: provider,
                              accountId: externalAccount.provider_user_id,
                              scope: approved_scopes,
                              userId: createdUser?.id,
                              createdAt: safeDateConversion(created_at),
                              updatedAt: safeDateConversion(updated_at),
                              password: password_digest,
                          }
                      })
                  } else {
                      await ctx.adapter.create({
                          model: "account",
                          data: {
                              id,
                              providerId: provider.replace("oauth_", ""),
                              accountId: externalAccount.provider_user_id,
                              scope: approved_scopes,
                              userId: createdUser?.id,
                              createdAt: safeDateConversion(created_at),
                              updatedAt: safeDateConversion(updated_at),
                          },
                          forceAllowId: true
                      })
                  }
              }
          }
    
          //two factor
          if (isTwoFactorEnabled) {
              await ctx.adapter.create({
                  model: "twoFactor",
                  data: {
                      userId: createdUser?.id,
                      secret: totp_secret,
                      backupCodes: await generateBackupCodes(totp_secret)
                  }
              })
          }
      }
    }
    
    migrateFromClerk()
      .then(() => {
          console.log('Migration completed');
          process.exit(0);
      })
      .catch((error) => {
          console.error('Migration failed:', error);
          process.exit(1);
      });</content>
</page>

<page>
  <title>Creem | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/creem</url>
  <content>[Creem](https://creem.io/) is a financial OS that enables teams and individuals selling software globally to split revenue and collaborate on financial workflows without any tax compliance headaches. This plugin integrates Creem with Better Auth, bringing payment processing and subscription management directly into your authentication layer.

[

### Get support on Creem Discord or in our in-app live-chat

This plugin is maintained by the Creem team.  
Need help? Reach out to our team anytime on Discord.

](https://discord.gg/q3GKZs92Av)

*   **Database Persistence** - Automatically synchronize customer and subscription data with your database
*   **Access Management** - Automatically grant or revoke access to users based on their subscription status
*   **Customer Synchronization** - Synchronize Creem customer IDs with your database users
*   **Checkout Integration** - Create checkout sessions either automatically for authenticated users or manually for unauthenticated users
*   **Customer Portal** - Enable users to manage subscriptions, view invoices, and update payment methods
*   **Subscription Management** - Cancel, retrieve, and track subscription details for authenticated users or manually for unauthenticated users
*   **Transaction History** - Search and filter transaction records for authenticated users or manually for unauthenticated users
*   **Webhook Processing** - Handle Creem webhooks securely with signature verification
*   **Flexible Architecture** - Use Better Auth endpoints or direct server-side functions
*   **Trial Abuse Prevention** - Users can only get one trial per account across all plans (when using database mode)

### [Install the plugin](#install-the-plugin)

If you're using a separate client and server setup, make sure to install the plugin in both parts of your project.

### [Get your API Key](#get-your-api-key)

Get your Creem API Key from the [Creem dashboard](https://creem.io/dashboard/developers), under the 'Developers' menu and add it to your environment variables:

    # .env
    CREEM_API_KEY=your_api_key_here

Test Mode and Production have different API keys. Make sure you're using the correct one for your environment.

### [Server Configuration](#server-configuration)

Configure Better Auth with the Creem plugin:

    // lib/auth.ts
    import { betterAuth } from "better-auth";
    import { creem } from "@creem_io/better-auth";
    
    export const auth = betterAuth({
      database: {
        // your database config
      },
      plugins: [
        creem({
          apiKey: process.env.CREEM_API_KEY!,
          webhookSecret: process.env.CREEM_WEBHOOK_SECRET, // Optional, webhooks are automatically enabled when passing a signing secret
          testMode: true, // Optional, use test mode for development
          defaultSuccessUrl: "/success", // Optional, redirect to this URL after successful payments
          persistSubscriptions: true, // Optional, enable database persistence (default: true)
        }),
      ],
    });

### [Client Configuration](#client-configuration)

### [Standard Setup](#standard-setup)

    // lib/auth-client.ts
    import { createAuthClient } from "better-auth/react";
    import { creemClient } from "@creem_io/better-auth/client";
    
    export const authClient = createAuthClient({
      baseURL: process.env.NEXT_PUBLIC_APP_URL,
      plugins: [creemClient()],
    });

### [Enhanced TypeScript Support (React-Only)](#enhanced-typescript-support-react-only)

For improved TypeScript IntelliSense and autocomplete:

    // lib/auth-client.ts
    import { createCreemAuthClient } from "@creem_io/better-auth/create-creem-auth-client";
    import { creemClient } from "@creem_io/better-auth/client";
    
    export const authClient = createCreemAuthClient({
      baseURL: process.env.NEXT_PUBLIC_APP_URL,
      plugins: [creemClient()],
    });

The `createCreemAuthClient` wrapper provides enhanced TypeScript support and cleaner parameter types. It's optimized for use with the Creem plugin.

### [Database Migration](#database-migration)

If you're using database persistence (`persistSubscriptions: true`), generate and run the database schema:

### [Webhook Setup](#webhook-setup)

### [Create Webhook Endpoint](#create-webhook-endpoint)

In your [Creem dashboard](https://creem.io/dashboard/developers/webhooks), create a webhook endpoint pointing to your local or production server pointing to:

    https://your-domain.com/api/auth/creem/webhook

(`/api/auth` is the default Better Auth server path)

Check step 3 if local development.

### [Configure Webhook Secret](#configure-webhook-secret)

Copy the webhook signing secret from Creem and add it to your environment:

    CREEM_WEBHOOK_SECRET=your_webhook_secret_here

Update your server configuration:

    creem({
      apiKey: process.env.CREEM_API_KEY!,
      webhookSecret: process.env.CREEM_WEBHOOK_SECRET,
      testMode: true,
    })

### [Local Development (Optional)](#local-development-optional)

For local testing, use a tool like [ngrok](https://ngrok.com/) to expose your local server:

    ngrok http 3000

Add the ngrok URL to your Creem webhook settings.

When `persistSubscriptions: true`, the plugin creates the following schema:

### [Subscription Table](#subscription-table)

| Field | Type | Description |
| --- | --- | --- |
| `id` | string | Primary key |
| `productId` | string | Creem product ID |
| `referenceId` | string | Your user/organization ID |
| `creemCustomerId` | string | Creem customer ID |
| `creemSubscriptionId` | string | Creem subscription ID |
| `creemOrderId` | string | Creem order ID |
| `status` | string | Subscription status |
| `periodStart` | date | Billing period start date |
| `periodEnd` | date | Billing period end date |
| `cancelAtPeriodEnd` | boolean | Whether subscription will cancel |

### [User Table Extension](#user-table-extension)

| Field | Type | Description |
| --- | --- | --- |
| `creemCustomerId` | string | Links user to Creem customer |

### [Checkout](#checkout)

Create a checkout session to process payments:

    "use client";
    
    import { authClient } from "@/lib/auth-client";
    
    export function SubscribeButton({ productId }: { productId: string }) {
      const handleCheckout = async () => {
        const { data, error } = await authClient.creem.createCheckout({
          productId,
          successUrl: "/dashboard",
          discountCode: "LAUNCH50", // Optional
          metadata: { planType: "pro" }, // Optional
        });
    
        if (data?.url) {
          window.location.href = data.url;
        }
      };
    
      return <button onClick={handleCheckout}>Subscribe Now</button>;
    }

#### [Checkout Options](#checkout-options)

*   `productId` (required) - The Creem product ID
*   `units` - Number of units (default: 1)
*   `successUrl` - Redirect URL after successful payment
*   `discountCode` - Discount code to apply
*   `customer` - Customer information (auto-populated from session)
*   `metadata` - Additional metadata (auto-includes user ID as `referenceId`)
*   `requestId` - Idempotency key for duplicate prevention

### [Customer Portal](#customer-portal)

Redirect users to manage their subscriptions:

    const handlePortal = async () => {
      // No need to redirect, the portal will be opened in the same tab
      const { data, error } = await authClient.creem.createPortal();
    };

### [Subscription Management](#subscription-management)

### [Cancel Subscription](#cancel-subscription)

When database persistence is enabled, the subscription is found automatically for the authenticated user:

    const handleCancel = async () => {
      const { data, error } = await authClient.creem.cancelSubscription();
    
      if (data?.success) {
        console.log(data.message);
      }
    };

If database persistence is disabled, provide the subscription ID:

    const { data } = await authClient.creem.cancelSubscription({
      id: "sub_123456",
    });

### [Retrieve Subscription](#retrieve-subscription)

Get subscription details for the authenticated user:

    const getSubscription = async () => {
      const { data } = await authClient.creem.retrieveSubscription();
    
      if (data) {
        console.log(`Status: ${data.status}`);
        console.log(`Product: ${data.product.name}`);
        console.log(`Price: ${data.product.price} ${data.product.currency}`);
      }
    };

### [Check Access](#check-access)

Verify if the user has an active subscription (requires database mode):

    const { data } = await authClient.creem.hasAccessGranted();
    
    if (data?.hasAccess) {
      // User has active subscription access
      console.log(`Expires: ${data.expiresAt}`);
    }

This function checks if the user has access for the current billing period. For example, if a user purchases a yearly plan and cancels after one month, they still have access until the year ends.

### [Transaction History](#transaction-history)

Search transaction records for the authenticated user:

    const { data } = await authClient.creem.searchTransactions({
      productId: "prod_xyz789", // Optional filter
      pageNumber: 1,
      pageSize: 50,
    });
    
    if (data?.transactions) {
      data.transactions.forEach((tx) => {
        console.log(`${tx.type}: ${tx.amount} ${tx.currency}`);
      });
    }

The plugin provides flexible webhook handling with both granular event handlers and high-level access control handlers.

### [High-Level Access Control Handlers (Recommended)](#high-level-access-control-handlers-recommended)

These handlers provide the simplest and most powerful way to manage user access. They automatically handle all payment scenarios and subscription states, so you don't need to manage individual subscription events.

**Database Persistence Required:** These handlers require the database persistence option to be enabled in your plugin configuration.

| Handler Name | Data Parameter Type | Description |
| --- | --- | --- |
| **`onGrantAccess`** | **`GrantAccessContext`** | **Called when a user should be granted access.** Handles successful payments, active subscriptions, and trial periods. Use this to enable features, add user to groups, or update permissions. |
| **`onRevokeAccess`** | **`RevokeAccessContext`** | **Called when a user's access should be revoked.** Handles cancellations, expirations, refunds, and failed payments. Use this to disable features, remove from groups, or revoke permissions. |

**Why use these handlers?**

*   Single source of truth for access control
*   Handles all payment scenarios automatically
*   Reduces code complexity and potential bugs
*   Works for both one-time purchases and subscriptions
*   Takes current billing period and access expiration dates into consideration

    // lib/auth.ts
    import { betterAuth } from "better-auth";
    import { creem } from "@creem_io/better-auth";
    
    export const auth = betterAuth({
      database: {
        // your database config
      },
      plugins:[ 
        creem({
          apiKey: process.env.CREEM_API_KEY!,
          webhookSecret: process.env.CREEM_WEBHOOK_SECRET!,
    
          onGrantAccess: async ({ reason, product, customer, metadata }) => {
            const userId = metadata?.referenceId as string;
    
            // Update your database specific logic
            await db.user.update({
              where: { id: userId },
              data: { 
                hasAccess: true, 
                subscriptionTier: product.name,
                accessReason: reason 
              },
            });
    
            console.log(`Granted ${reason} access to ${customer.email}`);
          },
    
          onRevokeAccess: async ({ reason, product, customer, metadata }) => {
            const userId = metadata?.referenceId as string;
    
            // Update your database specific logic
            await db.user.update({
              where: { id: userId },
              data: { 
                hasAccess: false, 
                revokeReason: reason 
              },
            });
    
            console.log(`Revoked access (${reason}) from ${customer.email}`);
          },
        }),
      ],
    })

### [Grant Access Reasons](#grant-access-reasons)

*   `subscription_active` - Subscription is active
*   `subscription_trialing` - Subscription is in trial period
*   `subscription_paid` - Subscription payment received

### [Revoke Access Reasons](#revoke-access-reasons)

*   `subscription_paused` - Subscription paused by user or admin
*   `subscription_expired` - Subscription expired without renewal
*   `subscription_period_end` - Current subscription period ended without renewal

* * *

### [Granular Event Handlers](#granular-event-handlers)

For advanced use cases where you need fine-grained control over specific events, use these handlers:

| Handler Name | Data Parameter Type | Description |
| --- | --- | --- |
| `onCheckoutCompleted` | `FlatCheckoutCompleted` | Called when a checkout is completed successfully. |
| `onRefundCreated` | `FlatRefundCreated` | Triggered when a refund is issued for a payment. |
| `onDisputeCreated` | `FlatDisputeCreated` | Invoked when a payment dispute/chargeback is created. |
| `onSubscriptionActive` | `FlatSubscriptionEvent` | Fired when a subscription becomes active. |
| `onSubscriptionTrialing` | `FlatSubscriptionEvent` | Subscription enters a trial period. |
| `onSubscriptionCanceled` | `FlatSubscriptionEvent` | Called when a subscription is canceled. |
| `onSubscriptionPaid` | `FlatSubscriptionEvent` | Subscription payment is received. |
| `onSubscriptionExpired` | `FlatSubscriptionEvent` | Subscription has expired (no renewal/payment). |
| `onSubscriptionUnpaid` | `FlatSubscriptionEvent` | Payment for a subscription failed or remains unpaid. |
| `onSubscriptionUpdate` | `FlatSubscriptionEvent` | Subscription settings/details updated. |
| `onSubscriptionPastDue` | `FlatSubscriptionEvent` | Subscription payment is late or overdue. |
| `onSubscriptionPaused` | `FlatSubscriptionEvent` | Subscription has been paused (by user or admin). |

### [How to use a Webhook Handler](#how-to-use-a-webhook-handler)

Handle individual webhook events with all properties flattened for easy access:

    // lib/auth.ts
    import { betterAuth } from "better-auth";
    import { creem } from "@creem_io/better-auth";
    
    export const auth = betterAuth({
      database: {
        // your database config
      },
      plugins: [
        creem({
          apiKey: process.env.CREEM_API_KEY!,
          webhookSecret: process.env.CREEM_WEBHOOK_SECRET!,
    
          onCheckoutCompleted: async (data) => {
            const { customer, product, order, webhookEventType } = data;
            console.log(`${customer.email} purchased ${product.name}`);
            
            // Perfect for one-time payments
            await sendThankYouEmail(customer.email);
          },
    
          onSubscriptionActive: async (data) => {
            const { customer, product, status } = data;
            // Handle active subscription
          },
    
          onSubscriptionTrialing: async (data) => {
            // Handle trial period
          },
    
          onSubscriptionCanceled: async (data) => {
            // Handle cancellation
          },
    
          onSubscriptionExpired: async (data) => {
            // Handle expiration
          },
    
          onRefundCreated: async (data) => {
            // Handle refunds
          },
    
          onDisputeCreated: async (data) => {
            // Handle disputes
          },
        }),
      ],
    });

### [Custom Webhook Handler](#custom-webhook-handler)

Create your own webhook endpoint with signature verification:

    // app/api/webhooks/custom/route.ts
    import { validateWebhookSignature } from "@creem_io/better-auth/server";
    
    export async function POST(req: Request) {
      const payload = await req.text();
      const signature = req.headers.get("creem-signature");
    
      if (
        !validateWebhookSignature(
          payload,
          signature,
          process.env.CREEM_WEBHOOK_SECRET!
        )
      ) {
        return new Response("Invalid signature", { status: 401 });
      }
    
      const event = JSON.parse(payload);
      // Your custom webhook handling logic
    
      return Response.json({ received: true });
    }

Use these utilities directly in Server Components, Server Actions, or API routes without going through Better Auth endpoints.

### [Import Server Utilities](#import-server-utilities)

    import {
      createCheckout,
      createPortal,
      cancelSubscription,
      retrieveSubscription,
      searchTransactions,
      checkSubscriptionAccess,
      isActiveSubscription,
      formatCreemDate,
      getDaysUntilRenewal,
      validateWebhookSignature,
    } from "@creem_io/better-auth/server";

### [Server Component Example](#server-component-example)

    import { checkSubscriptionAccess } from "@creem_io/better-auth/server";
    import { auth } from "@/lib/auth";
    import { headers } from "next/headers";
    import { redirect } from "next/navigation";
    
    export default async function DashboardPage() {
      const session = await auth.api.getSession({ headers: await headers() });
    
      if (!session?.user) {
        redirect("/login");
      }
    
      const status = await checkSubscriptionAccess(
        {
          apiKey: process.env.CREEM_API_KEY!,
          testMode: true,
        },
        {
          database: auth.options.database,
          userId: session.user.id,
        }
      );
    
      if (!status.hasAccess) {
        redirect("/subscribe");
      }
    
      return (
        <div>
          <h1>Welcome to Dashboard</h1>
          <p>Subscription Status: {status.status}</p>
          {status.expiresAt && (
            <p>Renews: {status.expiresAt.toLocaleDateString()}</p>
          )}
        </div>
      );
    }

### [Server Action Example](#server-action-example)

    "use server";
    
    import { createCheckout } from "@creem_io/better-auth/server";
    import { auth } from "@/lib/auth";
    import { headers } from "next/headers";
    import { redirect } from "next/navigation";
    
    export async function startCheckout(productId: string) {
      const session = await auth.api.getSession({ headers: await headers() });
    
      if (!session?.user) {
        throw new Error("Not authenticated");
      }
    
      const { url } = await createCheckout(
        {
          apiKey: process.env.CREEM_API_KEY!,
          testMode: true,
        },
        {
          productId,
          customer: { email: session.user.email },
          successUrl: "/success",
          metadata: { userId: session.user.id },
        }
      );
    
      redirect(url);
    }

### [Middleware Example](#middleware-example)

Protect routes based on subscription status:

    import { checkSubscriptionAccess } from "@creem_io/better-auth/server";
    import { auth } from "@/lib/auth";
    import { NextRequest, NextResponse } from "next/server";
    
    export async function middleware(request: NextRequest) {
      const session = await auth.api.getSession({
        headers: request.headers,
      });
    
      if (!session?.user) {
        return NextResponse.redirect(new URL("/login", request.url));
      }
    
      const status = await checkSubscriptionAccess(
        {
          apiKey: process.env.CREEM_API_KEY!,
          testMode: true,
        },
        {
          database: auth.options.database,
          userId: session.user.id,
        }
      );
    
      if (!status.hasAccess) {
        return NextResponse.redirect(new URL("/subscribe", request.url));
      }
    
      return NextResponse.next();
    }
    
    export const config = {
      matcher: ["/dashboard/:path*"],
    };

### [Utility Functions](#utility-functions)

    import {
      isActiveSubscription,
      formatCreemDate,
      getDaysUntilRenewal,
    } from "@creem_io/better-auth/server";
    
    // Check if status grants access
    if (isActiveSubscription(subscription.status)) {
      // User has access
    }
    
    // Format Creem timestamps
    const renewalDate = formatCreemDate(subscription.next_billing_date);
    console.log(renewalDate.toLocaleDateString());
    
    // Calculate days until renewal
    const days = getDaysUntilRenewal(subscription.current_period_end_date);
    console.log(`Renews in ${days} days`);

### [Database Mode vs API Mode](#database-mode-vs-api-mode)

The plugin supports two operational modes:

### [Database Mode (Recommended)](#database-mode-recommended)

When `persistSubscriptions: true` (default), subscription data is stored in your database.

**Benefits:**

*   Fast access checks without API calls
*   Offline access to subscription data
*   Query subscriptions with SQL
*   Automatic synchronization via webhooks
*   Trial abuse prevention

**Usage:**

    creem({
      apiKey: process.env.CREEM_API_KEY!,
      persistSubscriptions: true, // Default
    })

### [API Mode](#api-mode)

When `persistSubscriptions: false`, all data comes directly from the Creem API.

**Benefits:**

*   No database schema required
*   Simpler initial setup

**Limitations:**

*   Requires API call for each access check
*   Some features require custom implementation
*   No built-in trial abuse prevention

**Usage:**

    creem({
      apiKey: process.env.CREEM_API_KEY!,
      persistSubscriptions: false,
    })

In API mode, functions like `checkSubscriptionAccess` and `hasAccessGranted` have limited functionality and may require custom implementation using the Creem SDK directly.

### [Server-Side Types](#server-side-types)

| Type Name | Description | Typical Usage |
| --- | --- | --- |
| `CreemOptions` | Configuration options for the Creem plugin, such as API keys and persistence settings. | Used to configure the plugin on the server. |
| `GrantAccessContext` | Context passed to custom access control hooks when granting access to a user. | Used in custom access logic. |
| `RevokeAccessContext` | Context passed to hooks when revoking user access due to subscription status changes. | Used in custom access logic. |
| `GrantAccessReason` | Enum or type describing reasons for granting access (e.g., payment received, trial activated). | Returned in access-related hooks and events. |
| `RevokeAccessReason` | Enum or type describing reasons for revoking access (e.g., canceled, payment failed). | Returned in access-related hooks and events. |
| `FlatCheckoutCompleted` | Event object type for webhook payload when a checkout completes successfully. | Used in webhook handlers and event listeners. |
| `FlatRefundCreated` | Event object type for webhook payload when a refund is created. | Used in webhook handlers and event listeners. |
| `FlatDisputeCreated` | Event object type for webhook payload when a dispute is created. | Used in webhook handlers and event listeners. |
| `FlatSubscriptionEvent` | Event object type for generic subscription events (created, updated, canceled, etc). | Used in webhook handlers and event listeners. |

### [Client-Side Types](#client-side-types)

| Type Name | Description |
| --- | --- |
| `CreateCheckoutInput` | Input parameters for creating a checkout session. |
| `CreateCheckoutResponse` | Response shape for a checkout session creation request. |
| `CheckoutCustomer` | Customer information type used in a checkout session. |
| `CreatePortalInput` | Input parameters for creating a customer portal session. |
| `CreatePortalResponse` | Response data for a request to create a customer portal. |
| `CancelSubscriptionInput` | Input parameters when cancelling a subscription. |
| `CancelSubscriptionResponse` | Response data for a subscription cancellation request. |
| `RetrieveSubscriptionInput` | Input for retrieving a specific subscription's details. |
| `SubscriptionData` | Subscription information structure as returned by the API. |
| `SearchTransactionsInput` | Filters and parameters for searching transactions. |
| `SearchTransactionsResponse` | Response structure for a transaction search query. |
| `TransactionData` | Data relating to individual transactions (e.g., payment, refund, etc). |
| `HasAccessGrantedResponse` | The shape of the response indicating whether a user has access based on subscription status/rules. |

When using database mode (`persistSubscriptions: true`), the plugin automatically prevents trial abuse. Users can only receive one trial across all subscription plans.

**Example Scenario:**

1.  User subscribes to "Starter" plan with 7-day trial
2.  User cancels subscription during the trial period
3.  User attempts to subscribe to "Premium" plan
4.  No trial is offered - user is charged immediately

This protection is automatic and requires no configuration. Trial eligibility is determined when the subscription is created and cannot be overridden.

### [Webhook Issues](#webhook-issues)

If webhooks aren't being processed correctly:

1.  Verify the webhook URL is correct in your Creem dashboard
2.  Check that the webhook signing secret matches
3.  Ensure all necessary events are selected in the Creem dashboard
4.  Review server logs for webhook processing errors
5.  Test webhook delivery using Creem's webhook testing tool

### [Subscription Status Issues](#subscription-status-issues)

If subscription statuses aren't updating:

1.  Confirm webhooks are being received and processed
2.  Verify `creemCustomerId` and `creemSubscriptionId` fields are populated
3.  Check that reference IDs match between your application and Creem
4.  Review webhook handler logs for errors

### [Database Mode Not Working](#database-mode-not-working)

If database persistence isn't functioning:

1.  Ensure `persistSubscriptions: true` is set (it's the default)
2.  Run migrations: `npx @better-auth/cli migrate`
3.  Verify database connection is working
4.  Check that schema tables were created successfully
5.  Review database adapter configuration

### [API Mode Limitations](#api-mode-limitations)

Some functionalities are only available in database mode or require extra parameters to be passed:

*   `checkSubscriptionAccess` requires passing the `userId` parameter
*   `getActiveSubscriptions` requires passing the `userId` parameter
*   No automatic trial abuse prevention
*   No access to `hasAccessGranted` client method

To use these features, either enable database mode or implement custom logic using the Creem SDK directly.

*   [Creem Documentation](https://docs.creem.io/)
*   [Creem Dashboard](https://creem.io/dashboard)
*   [Better Auth Documentation](https://better-auth.com/)
*   [Plugin GitHub Repository Additional Documentation](https://github.com/armitage-labs/creem-betterauth)

For issues or questions:

*   Open an issue on [GitHub](https://github.com/armitage-labs/creem-betterauth/issues)
*   Contact Creem support at [support@creem.io](mailto:support@creem.io)
*   Join our [Discord community](https://discord.gg/q3GKZs92Av) for real-time support and discussion.
*   Chat with us directly using the in-app live chat on the [Creem dashboard](https://creem.io/dashboard).</content>
</page>

<page>
  <title>Organization | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/organization</url>
  <content>Organizations simplifies user access and permissions management. Assign roles and permissions to streamline project management, team coordination, and partnerships.

### [Add the plugin to your **auth** config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { organization } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [ 
            organization() 
        ] 
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { organizationClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [ 
            organizationClient() 
        ] 
    })

Once you've installed the plugin, you can start using the organization plugin to manage your organization's members and teams. The client plugin will provide you with methods under the `organization` namespace, and the server `api` will provide you with the necessary endpoints to manage your organization and give you an easier way to call the functions on your own backend.

### [Create an organization](#create-an-organization)

    const metadata = { someKey: "someValue" };const { data, error } = await authClient.organization.create({    name: "My Organization", // required    slug: "my-org", // required    logo: "https://example.com/logo.png",    metadata,    keepCurrentActiveOrganization: false,});

| Prop | Description | Type |
| --- | --- | --- |
| `name` | 
The organization name.

 | `string` |
| `slug` | 

The organization slug.

 | `string` |
| `logo?` | 

The organization logo.

 | `string` |
| `metadata?` | 

The metadata of the organization.

 | `Record<string, any>` |
| `keepCurrentActiveOrganization?` | 

Whether to keep the current active organization active after creating a new one.

 | `boolean` |

#### [Restrict who can create an organization](#restrict-who-can-create-an-organization)

By default, any user can create an organization. To restrict this, set the `allowUserToCreateOrganization` option to a function that returns a boolean, or directly to `true` or `false`.

auth.ts

    import { betterAuth } from "better-auth";
    import { organization } from "better-auth/plugins";
    
    const auth = betterAuth({
      //...
      plugins: [
        organization({
          allowUserToCreateOrganization: async (user) => {
            const subscription = await getSubscription(user.id); 
            return subscription.plan === "pro"; 
          }, 
        }),
      ],
    });

#### [Check if organization slug is taken](#check-if-organization-slug-is-taken)

To check if an organization slug is taken or not you can use the `checkSlug` function provided by the client. The function takes an object with the following properties:

POST

/organization/check-slug

    const { data, error } = await authClient.organization.checkSlug({    slug: "my-org", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `slug` | 
The organization slug to check.

 | `string` |

### [Organization Hooks](#organization-hooks)

You can customize organization operations using hooks that run before and after various organization-related activities. Better Auth provides two ways to configure hooks:

1.  **Legacy organizationCreation hooks** (deprecated, use `organizationHooks` instead)
2.  **Modern organizationHooks** (recommended) - provides comprehensive control over all organization-related activities

#### [Organization Creation and Management Hooks](#organization-creation-and-management-hooks)

Control organization lifecycle operations:

auth.ts

    import { betterAuth } from "better-auth";
    import { organization } from "better-auth/plugins";
    
    export const auth = betterAuth({
      plugins: [
        organization({
          organizationHooks: {
            // Organization creation hooks
            beforeCreateOrganization: async ({ organization, user }) => {
              // Run custom logic before organization is created
              // Optionally modify the organization data
              return {
                data: {
                  ...organization,
                  metadata: {
                    customField: "value",
                  },
                },
              };
            },
    
            afterCreateOrganization: async ({ organization, member, user }) => {
              // Run custom logic after organization is created
              // e.g., create default resources, send notifications
              await setupDefaultResources(organization.id);
            },
    
            // Organization update hooks
            beforeUpdateOrganization: async ({ organization, user, member }) => {
              // Validate updates, apply business rules
              return {
                data: {
                  ...organization,
                  name: organization.name?.toLowerCase(),
                },
              };
            },
    
            afterUpdateOrganization: async ({ organization, user, member }) => {
              // Sync changes to external systems
              await syncOrganizationToExternalSystems(organization);
            },
          },
        }),
      ],
    });

The legacy `organizationCreation` hooks are still supported but deprecated. Use `organizationHooks.beforeCreateOrganization` and `organizationHooks.afterCreateOrganization` instead for new projects.

#### [Member Hooks](#member-hooks)

Control member operations within organizations:

auth.ts

    import { betterAuth } from "better-auth";
    import { organization } from "better-auth/plugins";
    
    export const auth = betterAuth({
      plugins: [
        organization({
          organizationHooks: {
            // Before a member is added to an organization
            beforeAddMember: async ({ member, user, organization }) => {
              // Custom validation or modification
              console.log(`Adding ${user.email} to ${organization.name}`);
    
              // Optionally modify member data
              return {
                data: {
                  ...member,
                  role: "custom-role", // Override the role
                },
              };
            },
    
            // After a member is added
            afterAddMember: async ({ member, user, organization }) => {
              // Send welcome email, create default resources, etc.
              await sendWelcomeEmail(user.email, organization.name);
            },
    
            // Before a member is removed
            beforeRemoveMember: async ({ member, user, organization }) => {
              // Cleanup user's resources, send notification, etc.
              await cleanupUserResources(user.id, organization.id);
            },
    
            // After a member is removed
            afterRemoveMember: async ({ member, user, organization }) => {
              await logMemberRemoval(user.id, organization.id);
            },
    
            // Before updating a member's role
            beforeUpdateMemberRole: async ({
              member,
              newRole,
              user,
              organization,
            }) => {
              // Validate role change permissions
              if (newRole === "owner" && !hasOwnerUpgradePermission(user)) {
                throw new Error("Cannot upgrade to owner role");
              }
    
              // Optionally modify the role
              return {
                data: {
                  role: newRole,
                },
              };
            },
    
            // After updating a member's role
            afterUpdateMemberRole: async ({
              member,
              previousRole,
              user,
              organization,
            }) => {
              await logRoleChange(user.id, previousRole, member.role);
            },
          },
        }),
      ],
    });

#### [Invitation Hooks](#invitation-hooks)

Control invitation lifecycle:

auth.ts

    export const auth = betterAuth({
      plugins: [
        organization({
          organizationHooks: {
            // Before creating an invitation
            beforeCreateInvitation: async ({
              invitation,
              inviter,
              organization,
            }) => {
              // Custom validation or expiration logic
              const customExpiration = new Date(
                Date.now() + 1000 * 60 * 60 * 24 * 7
              ); // 7 days
    
              return {
                data: {
                  ...invitation,
                  expiresAt: customExpiration,
                },
              };
            },
    
            // After creating an invitation
            afterCreateInvitation: async ({
              invitation,
              inviter,
              organization,
            }) => {
              // Send custom invitation email, track metrics, etc.
              await sendCustomInvitationEmail(invitation, organization);
            },
    
            // Before accepting an invitation
            beforeAcceptInvitation: async ({ invitation, user, organization }) => {
              // Additional validation before acceptance
              await validateUserEligibility(user, organization);
            },
    
            // After accepting an invitation
            afterAcceptInvitation: async ({
              invitation,
              member,
              user,
              organization,
            }) => {
              // Setup user account, assign default resources
              await setupNewMemberResources(user, organization);
            },
    
            // Before/after rejecting invitations
            beforeRejectInvitation: async ({ invitation, user, organization }) => {
              // Log rejection reason, send notification to inviter
            },
    
            afterRejectInvitation: async ({ invitation, user, organization }) => {
              await notifyInviterOfRejection(invitation.inviterId, user.email);
            },
    
            // Before/after cancelling invitations
            beforeCancelInvitation: async ({
              invitation,
              cancelledBy,
              organization,
            }) => {
              // Verify cancellation permissions
            },
    
            afterCancelInvitation: async ({
              invitation,
              cancelledBy,
              organization,
            }) => {
              await logInvitationCancellation(invitation.id, cancelledBy.id);
            },
          },
        }),
      ],
    });

#### [Team Hooks](#team-hooks)

Control team operations (when teams are enabled):

auth.ts

    export const auth = betterAuth({
      plugins: [
        organization({
          teams: { enabled: true },
          organizationHooks: {
            // Before creating a team
            beforeCreateTeam: async ({ team, user, organization }) => {
              // Validate team name, apply naming conventions
              return {
                data: {
                  ...team,
                  name: team.name.toLowerCase().replace(/\s+/g, "-"),
                },
              };
            },
    
            // After creating a team
            afterCreateTeam: async ({ team, user, organization }) => {
              // Create default team resources, channels, etc.
              await createDefaultTeamResources(team.id);
            },
    
            // Before updating a team
            beforeUpdateTeam: async ({ team, updates, user, organization }) => {
              // Validate updates, apply business rules
              return {
                data: {
                  ...updates,
                  name: updates.name?.toLowerCase(),
                },
              };
            },
    
            // After updating a team
            afterUpdateTeam: async ({ team, user, organization }) => {
              await syncTeamChangesToExternalSystems(team);
            },
    
            // Before deleting a team
            beforeDeleteTeam: async ({ team, user, organization }) => {
              // Backup team data, notify members
              await backupTeamData(team.id);
            },
    
            // After deleting a team
            afterDeleteTeam: async ({ team, user, organization }) => {
              await cleanupTeamResources(team.id);
            },
    
            // Team member operations
            beforeAddTeamMember: async ({
              teamMember,
              team,
              user,
              organization,
            }) => {
              // Validate team membership limits, permissions
              const memberCount = await getTeamMemberCount(team.id);
              if (memberCount >= 10) {
                throw new Error("Team is full");
              }
            },
    
            afterAddTeamMember: async ({
              teamMember,
              team,
              user,
              organization,
            }) => {
              await grantTeamAccess(user.id, team.id);
            },
    
            beforeRemoveTeamMember: async ({
              teamMember,
              team,
              user,
              organization,
            }) => {
              // Backup user's team-specific data
              await backupTeamMemberData(user.id, team.id);
            },
    
            afterRemoveTeamMember: async ({
              teamMember,
              team,
              user,
              organization,
            }) => {
              await revokeTeamAccess(user.id, team.id);
            },
          },
        }),
      ],
    });

#### [Hook Error Handling](#hook-error-handling)

All hooks support error handling. Throwing an error in a `before` hook will prevent the operation from proceeding:

auth.ts

    import { APIError } from "better-auth/api";
    
    export const auth = betterAuth({
      plugins: [
        organization({
          organizationHooks: {
            beforeAddMember: async ({ member, user, organization }) => {
              // Check if user has pending violations
              const violations = await checkUserViolations(user.id);
              if (violations.length > 0) {
                throw new APIError("BAD_REQUEST", {
                  message:
                    "User has pending violations and cannot join organizations",
                });
              }
            },
    
            beforeCreateTeam: async ({ team, user, organization }) => {
              // Validate team name uniqueness
              const existingTeam = await findTeamByName(team.name, organization.id);
              if (existingTeam) {
                throw new APIError("BAD_REQUEST", {
                  message: "Team name already exists in this organization",
                });
              }
            },
          },
        }),
      ],
    });

### [List User's Organizations](#list-users-organizations)

To list the organizations that a user is a member of, you can use `useListOrganizations` hook. It implements a reactive way to get the organizations that the user is a member of.

Or alternatively, you can call `organization.list` if you don't want to use a hook.

    const { data, error } = await authClient.organization.list();

### [Active Organization](#active-organization)

Active organization is the workspace the user is currently working on. By default when the user is signed in the active organization is set to `null`. You can set the active organization to the user session.

It's not always you want to persist the active organization in the session. You can manage the active organization in the client side only. For example, multiple tabs can have different active organizations.

#### [Set Active Organization](#set-active-organization)

You can set the active organization by calling the `organization.setActive` function. It'll set the active organization for the user session.

In some applications, you may want the ability to unset an active organization. In this case, you can call this endpoint with `organizationId` set to `null`.

POST

/organization/set-active

    const { data, error } = await authClient.organization.setActive({    organizationId: "org-id",    organizationSlug: "org-slug",});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId?` | 
The organization ID to set as active. It can be null to unset the active organization.

 | `string | null` |
| `organizationSlug?` | 

The organization slug to set as active. It can be null to unset the active organization if organizationId is not provided.

 | `string` |

To set active organization when a session is created you can use [database hooks](https://www.better-auth.com/docs/concepts/database#database-hooks).

auth.ts

    export const auth = betterAuth({
      databaseHooks: {
        session: {
          create: {
            before: async (session) => {
              const organization = await getActiveOrganization(session.userId);
              return {
                data: {
                  ...session,
                  activeOrganizationId: organization.id,
                },
              };
            },
          },
        },
      },
    });

#### [Use Active Organization](#use-active-organization)

To retrieve the active organization for the user, you can call the `useActiveOrganization` hook. It returns the active organization for the user. Whenever the active organization changes, the hook will re-evaluate and return the new active organization.

### [Get Full Organization](#get-full-organization)

To get the full details of an organization, you can use the `getFullOrganization` function. By default, if you don't pass any properties, it will use the active organization.

GET

/organization/get-full-organization

    const { data, error } = await authClient.organization.getFullOrganization({    query: {        organizationId: "org-id",        organizationSlug: "org-slug",        membersLimit: 100,    },});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId?` | 
The organization ID to get. By default, it will use the active organization.

 | `string` |
| `organizationSlug?` | 

The organization slug to get.

 | `string` |
| `membersLimit?` | 

The limit of members to get. By default, it uses the membershipLimit option which defaults to 100.

 | `number` |

### [Update Organization](#update-organization)

To update organization info, you can use `organization.update`

    const { data, error } = await authClient.organization.update({    data: { // required        name: "updated-name",        slug: "updated-slug",        logo: "new-logo.url",        metadata: { customerId: "test" },    },    organizationId: "org-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `data` | 
A partial list of data to update the organization.

 | `Object` |
| `data.name?` | 

The name of the organization.

 | `string` |
| `data.slug?` | 

The slug of the organization.

 | `string` |
| `data.logo?` | 

The logo of the organization.

 | `string` |
| `data.metadata?` | 

The metadata of the organization.

 | `Record<string, any> | null` |
| `organizationId?` | 

The organization ID. to update.

 | `string` |

### [Delete Organization](#delete-organization)

To remove user owned organization, you can use `organization.delete`

    const { data, error } = await authClient.organization.delete({    organizationId: "org-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId` | 
The organization ID to delete.

 | `string` |

If the user has the necessary permissions (by default: role is owner) in the specified organization, all members, invitations and organization information will be removed.

You can configure how organization deletion is handled through `organizationDeletion` option:

    const auth = betterAuth({
      plugins: [
        organization({
          disableOrganizationDeletion: true, //to disable it altogether
          organizationHooks: {
            beforeDeleteOrganization: async (data, request) => {
              // a callback to run before deleting org
            },
            afterDeleteOrganization: async (data, request) => {
              // a callback to run after deleting org
            },
          },
        }),
      ],
    });

To add a member to an organization, we first need to send an invitation to the user. The user will receive an email/sms with the invitation link. Once the user accepts the invitation, they will be added to the organization.

### [Setup Invitation Email](#setup-invitation-email)

For member invitation to work we first need to provide `sendInvitationEmail` to the `better-auth` instance. This function is responsible for sending the invitation email to the user.

You'll need to construct and send the invitation link to the user. The link should include the invitation ID, which will be used with the acceptInvitation function when the user clicks on it.

auth.ts

    import { betterAuth } from "better-auth";
    import { organization } from "better-auth/plugins";
    import { sendOrganizationInvitation } from "./email";
    export const auth = betterAuth({
      plugins: [
        organization({
          async sendInvitationEmail(data) {
            const inviteLink = `https://example.com/accept-invitation/${data.id}`;
            sendOrganizationInvitation({
              email: data.email,
              invitedByUsername: data.inviter.user.name,
              invitedByEmail: data.inviter.user.email,
              teamName: data.organization.name,
              inviteLink,
            });
          },
        }),
      ],
    });

### [Send Invitation](#send-invitation)

To invite users to an organization, you can use the `invite` function provided by the client. The `invite` function takes an object with the following properties:

POST

/organization/invite-member

    const { data, error } = await authClient.organization.inviteMember({    email: "example@gmail.com", // required    role: "member", // required    organizationId: "org-id",    resend: true,    teamId: "team-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `email` | 
The email address of the user to invite.

 | `string` |
| `role` | 

The role(s) to assign to the user. It can be `admin`, `member`, `owner`

 | `string | string[]` |
| `organizationId?` | 

The organization ID to invite the user to. Defaults to the active organization.

 | `string` |
| `resend?` | 

Resend the invitation email, if the user is already invited.

 | `boolean` |
| `teamId?` | 

The team ID to invite the user to.

 | `string` |

*   If the user is already a member of the organization, the invitation will be canceled. - If the user is already invited to the organization, unless `resend` is set to `true`, the invitation will not be sent again. - If `cancelPendingInvitationsOnReInvite` is set to `true`, the invitation will be canceled if the user is already invited to the organization and a new invitation is sent.

### [Accept Invitation](#accept-invitation)

When a user receives an invitation email, they can click on the invitation link to accept the invitation. The invitation link should include the invitation ID, which will be used to accept the invitation.

Make sure to call the `acceptInvitation` function after the user is logged in.

POST

/organization/accept-invitation

    const { data, error } = await authClient.organization.acceptInvitation({    invitationId: "invitation-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `invitationId` | 
The ID of the invitation to accept.

 | `string` |

#### [Email Verification Requirement](#email-verification-requirement)

If the `requireEmailVerificationOnInvitation` option is enabled in your organization configuration, users must verify their email address before they can accept invitations. This adds an extra security layer to ensure that only verified users can join your organization.

auth.ts

    import { betterAuth } from "better-auth";
    import { organization } from "better-auth/plugins";
    
    export const auth = betterAuth({
      plugins: [
        organization({
          requireEmailVerificationOnInvitation: true, 
          async sendInvitationEmail(data) {
            // ... your email sending logic
          },
        }),
      ],
    });

### [Cancel Invitation](#cancel-invitation)

If a user has sent out an invitation, you can use this method to cancel it.

If you're looking for how a user can reject an invitation, you can find that [here](#reject-invitation).

POST

/organization/cancel-invitation

    await authClient.organization.cancelInvitation({    invitationId: "invitation-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `invitationId` | 
The ID of the invitation to cancel.

 | `string` |

### [Reject Invitation](#reject-invitation)

If this user has received an invitation, but wants to decline it, this method will allow you to do so by rejecting it.

POST

/organization/reject-invitation

    await authClient.organization.rejectInvitation({    invitationId: "invitation-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `invitationId` | 
The ID of the invitation to reject.

 | `string` |

Like accepting invitations, rejecting invitations also requires email verification when the `requireEmailVerificationOnInvitation` option is enabled. Users with unverified emails will receive an error when attempting to reject invitations.

### [Get Invitation](#get-invitation)

To get an invitation you can use the `organization.getInvitation` function provided by the client. You need to provide the invitation id as a query parameter.

GET

/organization/get-invitation

    const { data, error } = await authClient.organization.getInvitation({    query: {        id: "invitation-id", // required    },});

| Prop | Description | Type |
| --- | --- | --- |
| `id` | 
The ID of the invitation to get.

 | `string` |

### [List Invitations](#list-invitations)

To list all invitations for a given organization you can use the `listInvitations` function provided by the client.

GET

/organization/list-invitations

    const { data, error } = await authClient.organization.listInvitations({    query: {        organizationId: "organization-id",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId?` | 
An optional ID of the organization to list invitations for. If not provided, will default to the user's active organization.

 | `string` |

### [List user invitations](#list-user-invitations)

To list all invitations for a given user you can use the `listUserInvitations` function provided by the client.

auth-client.ts

    const invitations = await authClient.organization.listUserInvitations();

On the server, you can pass the user ID as a query parameter.

api.ts

    const invitations = await auth.api.listUserInvitations({
      query: {
        email: "user@example.com",
      },
    });

The `email` query parameter is only available on the server to query for invitations for a specific user.

### [List Members](#list-members)

To list all members of an organization you can use the `listMembers` function.

GET

/organization/list-members

    const { data, error } = await authClient.organization.listMembers({    query: {        organizationId: "organization-id",        limit: 100,        offset: 0,        sortBy: "createdAt",        sortDirection: "desc",        filterField: "createdAt",        filterOperator: "eq",        filterValue: "value",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId?` | 
An optional organization ID to list members for. If not provided, will default to the user's active organization.

 | `string` |
| `limit?` | 

The limit of members to return.

 | `number` |
| `offset?` | 

The offset to start from.

 | `number` |
| `sortBy?` | 

The field to sort by.

 | `string` |
| `sortDirection?` | 

The direction to sort by.

 | `"asc" | "desc"` |
| `filterField?` | 

The field to filter by.

 | `string` |
| `filterOperator?` | 

The operator to filter by.

 | `"eq" | "ne" | "gt" | "gte" | "lt" | "lte" | "in" | "nin" | "contains"` |
| `filterValue?` | 

The value to filter by.

 | `string` |

### [Remove Member](#remove-member)

To remove you can use `organization.removeMember`

POST

/organization/remove-member

    const { data, error } = await authClient.organization.removeMember({    memberIdOrEmail: "user@example.com", // required    organizationId: "org-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `memberIdOrEmail` | 
The ID or email of the member to remove.

 | `string` |
| `organizationId?` | 

The ID of the organization to remove the member from. If not provided, the active organization will be used.

 | `string` |

### [Update Member Role](#update-member-role)

To update the role of a member in an organization, you can use the `organization.updateMemberRole`. If the user has the permission to update the role of the member, the role will be updated.

POST

/organization/update-member-role

    await authClient.organization.updateMemberRole({    role: ["admin", "sale"], // required    memberId: "member-id", // required    organizationId: "organization-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `role` | 
The new role to be applied. This can be a string or array of strings representing the roles.

 | `string | string[]` |
| `memberId` | 

The member id to apply the role update to.

 | `string` |
| `organizationId?` | 

An optional organization ID which the member is a part of to apply the role update. If not provided, you must provide session headers to get the active organization.

 | `string` |

### [Get Active Member](#get-active-member)

To get the current member of the active organization you can use the `organization.getActiveMember` function. This function will return the user's member details in their active organization.

GET

/organization/get-active-member

    const { data: member, error } = await authClient.organization.getActiveMember();

### [Get Active Member Role](#get-active-member-role)

To get the current role member of the active organization you can use the `organization.getActiveMemberRole` function. This function will return the user's member role in their active organization.

GET

/organization/get-active-member-role

    const { data: { role }, error } = await authClient.organization.getActiveMemberRole();

### [Add Member](#add-member)

If you want to add a member directly to an organization without sending an invitation, you can use the `addMember` function which can only be invoked on the server.

POST

/organization/add-member

    const data = await auth.api.addMember({    body: {        userId: "user-id",        role: ["admin", "sale"], // required        organizationId: "org-id",        teamId: "team-id",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `userId?` | 
The user ID which represents the user to be added as a member. If `null` is provided, then it's expected to provide session headers.

 | `string | null` |
| `role` | 

The role(s) to assign to the new member.

 | `string | string[]` |
| `organizationId?` | 

An optional organization ID to pass. If not provided, will default to the user's active organization.

 | `string` |
| `teamId?` | 

An optional team ID to add the member to.

 | `string` |

### [Leave Organization](#leave-organization)

To leave organization you can use `organization.leave` function. This function will remove the current user from the organization.

    await authClient.organization.leave({    organizationId: "organization-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId` | 
The organization ID for the member to leave.

 | `string` |

The organization plugin provides a very flexible access control system. You can control the access of the user based on the role they have in the organization. You can define your own set of permissions based on the role of the user.

### [Roles](#roles)

By default, there are three roles in the organization:

`owner`: The user who created the organization by default. The owner has full control over the organization and can perform any action.

`admin`: Users with the admin role have full control over the organization except for deleting the organization or changing the owner.

`member`: Users with the member role have limited control over the organization. They can create projects, invite users, and manage projects they have created.

A user can have multiple roles. Multiple roles are stored as string separated by comma (",").

### [Permissions](#permissions)

By default, there are three resources, and these have two to three actions.

**organization**:

`update` `delete`

**member**:

`create` `update` `delete`

**invitation**:

`create` `cancel`

The owner has full control over all the resources and actions. The admin has full control over all the resources except for deleting the organization or changing the owner. The member has no control over any of those actions other than reading the data.

### [Custom Permissions](#custom-permissions)

The plugin provides an easy way to define your own set of permissions for each role.

#### [Create Access Control](#create-access-control)

You first need to create access controller by calling `createAccessControl` function and passing the statement object. The statement object should have the resource name as the key and the array of actions as the value.

permissions.ts

    import { createAccessControl } from "better-auth/plugins/access";
    
    /**
     * make sure to use `as const` so typescript can infer the type correctly
     */
    const statement = { 
        project: ["create", "share", "update", "delete"], 
    } as const; 
    
    const ac = createAccessControl(statement); 

#### [Create Roles](#create-roles)

Once you have created the access controller you can create roles with the permissions you have defined.

permissions.ts

    import { createAccessControl } from "better-auth/plugins/access";
    
    const statement = {
        project: ["create", "share", "update", "delete"],
    } as const;
    
    const ac = createAccessControl(statement);
    
    const member = ac.newRole({ 
        project: ["create"], 
    }); 
    
    const admin = ac.newRole({ 
        project: ["create", "update"], 
    }); 
    
    const owner = ac.newRole({ 
        project: ["create", "update", "delete"], 
    }); 
    
    const myCustomRole = ac.newRole({ 
        project: ["create", "update", "delete"], 
        organization: ["update"], 
    }); 

When you create custom roles for existing roles, the predefined permissions for those roles will be overridden. To add the existing permissions to the custom role, you need to import `defaultStatements` and merge it with your new statement, plus merge the roles' permissions set with the default roles.

permissions.ts

    import { createAccessControl } from "better-auth/plugins/access";
    import { defaultStatements, adminAc } from 'better-auth/plugins/organization/access'
    
    const statement = {
        ...defaultStatements, 
        project: ["create", "share", "update", "delete"],
    } as const;
    
    const ac = createAccessControl(statement);
    
    const admin = ac.newRole({
        project: ["create", "update"],
        ...adminAc.statements, 
    });

#### [Pass Roles to the Plugin](#pass-roles-to-the-plugin)

Once you have created the roles you can pass them to the organization plugin both on the client and the server.

auth.ts

    import { betterAuth } from "better-auth"
    import { organization } from "better-auth/plugins"
    import { ac, owner, admin, member } from "@/auth/permissions"
    
    export const auth = betterAuth({
        plugins: [
            organization({
                ac,
                roles: {
                    owner,
                    admin,
                    member,
                    myCustomRole
                }
            }),
        ],
    });

You also need to pass the access controller and the roles to the client plugin.

auth-client

    import { createAuthClient } from "better-auth/client"
    import { organizationClient } from "better-auth/client/plugins"
    import { ac, owner, admin, member, myCustomRole } from "@/auth/permissions"
    
    export const authClient = createAuthClient({
        plugins: [
            organizationClient({
                ac,
                roles: {
                    owner,
                    admin,
                    member,
                    myCustomRole
                }
            })
      ]
    })

### [Access Control Usage](#access-control-usage)

**Has Permission**:

You can use the `hasPermission` action provided by the `api` to check the permission of the user.

api.ts

    import { auth } from "@/auth";
    
    await auth.api.hasPermission({
      headers: await headers(),
      body: {
        permissions: {
          project: ["create"], // This must match the structure in your access control
        },
      },
    });
    
    // You can also check multiple resource permissions at the same time
    await auth.api.hasPermission({
      headers: await headers(),
      body: {
        permissions: {
          project: ["create"], // This must match the structure in your access control
          sale: ["create"],
        },
      },
    });

If you want to check the permission of the user on the client from the server you can use the `hasPermission` function provided by the client.

auth-client.ts

    const canCreateProject = await authClient.organization.hasPermission({
      permissions: {
        project: ["create"],
      },
    });
    
    // You can also check multiple resource permissions at the same time
    const canCreateProjectAndCreateSale =
      await authClient.organization.hasPermission({
        permissions: {
          project: ["create"],
          sale: ["create"],
        },
      });

**Check Role Permission**:

Once you have defined the roles and permissions to avoid checking the permission from the server you can use the `checkRolePermission` function provided by the client.

auth-client.ts

    const canCreateProject = authClient.organization.checkRolePermission({
      permissions: {
        organization: ["delete"],
      },
      role: "admin",
    });
    
    // You can also check multiple resource permissions at the same time
    const canCreateProjectAndCreateSale =
      authClient.organization.checkRolePermission({
        permissions: {
          organization: ["delete"],
          member: ["delete"],
        },
        role: "admin",
      });

This will not include any dynamic roles as everything is ran synchronously on the client side. Please use the [hasPermission](#access-control-usage) APIs to include checks for any dynamic roles & permissions.

* * *

Dynamic access control allows you to create roles at runtime for organizations. This is achieved by storing the created roles and permissions associated with an organization in a database table.

### [Enabling Dynamic Access Control](#enabling-dynamic-access-control)

To enable dynamic access control, pass the `dynamicAccessControl` configuration option with `enabled` set to `true` to both server and client plugins.

Ensure you have pre-defined an `ac` instance on the server auth plugin. This is important as this is how we can infer the permissions that are available for use.

auth.ts

    import { betterAuth } from "better-auth";
    import { organization } from "better-auth/plugins";
    import { ac } from "@/auth/permissions";
    
    export const auth = betterAuth({
        plugins: [ 
            organization({ 
                ac, // Must be defined in order for dynamic access control to work
                dynamicAccessControl: { 
                  enabled: true, 
                }, 
            }) 
        ] 
    })

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { organizationClient } from "better-auth/client/plugins";
    
    export const authClient = createAuthClient({
        plugins: [ 
            organizationClient({ 
                dynamicAccessControl: { 
                  enabled: true, 
                }, 
            }) 
        ] 
    })

This will require you to run migrations to add the new `organizationRole` table to the database.

The `authClient.organization.checkRolePermission` function will not include any dynamic roles as everything is ran synchronously on the client side. Please use the [hasPermission](#access-control-usage) APIs to include checks for any dynamic roles.

### [Creating a role](#creating-a-role)

To create a new role for an organization at runtime, you can use the `createRole` function.

Only users with roles which contain the `ac` resource with the `create` permission can create a new role. By default, only the `admin` and `owner` roles have this permission. You also cannot add permissions that your current role in that organization can't already access.

TIP: You can validate role names by using the `dynamicAccessControl.validateRoleName` option in the organization plugin config. Learn more [here](#validaterolename).

POST

/organization/create-role

    // To use custom resources or permissions,// make sure they are defined in the `ac` instance of your organization config.const permission = {  project: ["create", "update", "delete"]}await authClient.organization.createRole({    role: "my-unique-role", // required    permission: permission,    organizationId: "organization-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `role` | 
A unique name of the role to create.

 | `string` |
| `permission?` | 

The permissions to assign to the role.

 | `Record<string, string[]>` |
| `organizationId?` | 

The organization ID which the role will be created in. Defaults to the active organization.

 | `string` |

Now you can freely call [`updateMemberRole`](#updating-a-member-role) to update the role of a member with your newly created role!

### [Deleting a role](#deleting-a-role)

To delete a role, you can use the `deleteRole` function, then provide either a `roleName` or `roleId` parameter along with the `organizationId` parameter.

POST

/organization/delete-role

    await authClient.organization.deleteRole({    roleName: "my-role",    roleId: "role-id",    organizationId: "organization-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `roleName?` | 
The name of the role to delete. Alternatively, you can pass a `roleId` parameter instead.

 | `string` |
| `roleId?` | 

The id of the role to delete. Alternatively, you can pass a `roleName` parameter instead.

 | `string` |
| `organizationId?` | 

The organization ID which the role will be deleted in. Defaults to the active organization.

 | `string` |

### [Listing roles](#listing-roles)

To list roles, you can use the `listOrgRoles` function. This requires the `ac` resource with the `read` permission for the member to be able to list roles.

GET

/organization/list-roles

    const { data: roles, error } = await authClient.organization.listRoles({    query: {        organizationId: "organization-id",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId?` | 
The organization ID which the roles are under to list. Defaults to the user's active organization.

 | `string` |

### [Getting a specific role](#getting-a-specific-role)

To get a specific role, you can use the `getOrgRole` function and pass either a `roleName` or `roleId` parameter. This requires the `ac` resource with the `read` permission for the member to be able to get a role.

GET

/organization/get-role

    const { data: role, error } = await authClient.organization.getRole({    query: {        roleName: "my-role",        roleId: "role-id",        organizationId: "organization-id",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `roleName?` | 
The name of the role to get. Alternatively, you can pass a `roleId` parameter instead.

 | `string` |
| `roleId?` | 

The id of the role to get. Alternatively, you can pass a `roleName` parameter instead.

 | `string` |
| `organizationId?` | 

The organization ID which the role will be deleted in. Defaults to the active organization.

 | `string` |

### [Updating a role](#updating-a-role)

To update a role, you can use the `updateOrgRole` function and pass either a `roleName` or `roleId` parameter.

POST

/organization/update-role

    const { data: updatedRole, error } = await authClient.organization.updateRole({    roleName: "my-role",    roleId: "role-id",    organizationId: "organization-id",    data: { // required        permission: { project: ["create", "update", "delete"] },        roleName: "my-new-role",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `roleName?` | 
The name of the role to update. Alternatively, you can pass a `roleId` parameter instead.

 | `string` |
| `roleId?` | 

The id of the role to update. Alternatively, you can pass a `roleName` parameter instead.

 | `string` |
| `organizationId?` | 

The organization ID which the role will be updated in. Defaults to the active organization.

 | `string` |
| `data` | 

The data which will be updated

 | `Object` |
| `data.permission?` | 

Optionally update the permissions of the role.

 | `Record<string, string[]>` |
| `data.roleName?` | 

Optionally update the name of the role.

 | `string` |

### [Configuration Options](#configuration-options)

Below is a list of options that can be passed to the `dynamicAccessControl` object.

#### [`enabled`](#enabled)

This option is used to enable or disable dynamic access control. By default, it is disabled.

    organization({
      dynamicAccessControl: {
        enabled: true
      }
    })

#### [`maximumRolesPerOrganization`](#maximumrolesperorganization)

This option is used to limit the number of roles that can be created for an organization.

By default, the maximum number of roles that can be created for an organization is infinite.

    organization({
      dynamicAccessControl: {
        maximumRolesPerOrganization: 10
      }
    })

You can also pass a function that returns a number.

    organization({
      dynamicAccessControl: {
        maximumRolesPerOrganization: async (organizationId) => { 
          const organization = await getOrganization(organizationId); 
          return organization.plan === "pro" ? 100 : 10; 
        } 
      }
    })

### [Additional Fields](#additional-fields)

To add additional fields to the `organizationRole` table, you can pass the `additionalFields` configuration option to the `organization` plugin.

    organization({
      schema: {
        organizationRole: {
          additionalFields: {
            // Role colors!
            color: {
              type: "string",
              defaultValue: "#ffffff",
            },
            //... other fields
          },
        },
      },
    })

Then, if you don't already use `inferOrgAdditionalFields` to infer the additional fields, you can use it to infer the additional fields.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { organizationClient, inferOrgAdditionalFields } from "better-auth/client/plugins"
    import type { auth } from "./auth"
    
    export const authClient = createAuthClient({
        plugins: [
            organizationClient({
                schema: inferOrgAdditionalFields<typeof auth>()
            })
        ]
    })

Otherwise, you can pass the schema values directly, the same way you do on the org plugin in the server.

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { organizationClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [
            organizationClient({
                schema: {
                    organizationRole: {
                        additionalFields: {
                            color: {
                                type: "string",
                                defaultValue: "#ffffff",
                            }
                        }
                    }
                }
            })
        ]
    })

* * *

Teams allow you to group members within an organization. The teams feature provides additional organization structure and can be used to manage permissions at a more granular level.

### [Enabling Teams](#enabling-teams)

To enable teams, pass the `teams` configuration option to both server and client plugins:

auth.ts

    import { betterAuth } from "better-auth";
    import { organization } from "better-auth/plugins";
    
    export const auth = betterAuth({
      plugins: [
        organization({
          teams: {
            enabled: true,
            maximumTeams: 10, // Optional: limit teams per organization
            allowRemovingAllTeams: false, // Optional: prevent removing the last team
          },
        }),
      ],
    });

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import { organizationClient } from "better-auth/client/plugins";
    
    export const authClient = createAuthClient({
      plugins: [
        organizationClient({
          teams: {
            enabled: true,
          },
        }),
      ],
    });

### [Managing Teams](#managing-teams)

#### [Create Team](#create-team)

Create a new team within an organization:

POST

/organization/create-team

    const { data, error } = await authClient.organization.createTeam({    name: "my-team", // required    organizationId: "organization-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `name` | 
The name of the team.

 | `string` |
| `organizationId?` | 

The organization ID which the team will be created in. Defaults to the active organization.

 | `string` |

#### [List Teams](#list-teams)

Get all teams in an organization:

GET

/organization/list-teams

    const { data, error } = await authClient.organization.listTeams({    query: {        organizationId: "organization-id",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `organizationId?` | 
The organization ID which the teams are under to list. Defaults to the user's active organization.

 | `string` |

#### [Update Team](#update-team)

Update a team's details:

POST

/organization/update-team

    const { data, error } = await authClient.organization.updateTeam({    teamId: "team-id", // required    data: { // required        name: "My new team name",        organizationId: "My new organization ID for this team",        createdAt: new Date(),        updatedAt: new Date(),    },});

| Prop | Description | Type |
| --- | --- | --- |
| `teamId` | 
The ID of the team to be updated.

 | `string` |
| `data` | 

A partial object containing options for you to update.

 | `Object` |
| `data.name?` | 

The name of the team to be updated.

 | `string` |
| `data.organizationId?` | 

The organization ID which the team falls under.

 | `string` |
| `data.createdAt?` | 

The timestamp of when the team was created.

 | `Date` |
| `data.updatedAt?` | 

The timestamp of when the team was last updated.

 | `Date` |

#### [Remove Team](#remove-team)

Delete a team from an organization:

POST

/organization/remove-team

    const { data, error } = await authClient.organization.removeTeam({    teamId: "team-id", // required    organizationId: "organization-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `teamId` | 
The team ID of the team to remove.

 | `string` |
| `organizationId?` | 

The organization ID which the team falls under. If not provided, it will default to the user's active organization.

 | `string` |

#### [Set Active Team](#set-active-team)

Sets the given team as the current active team. If `teamId` is `null` the current active team is unset.

POST

/organization/set-active-team

    const { data, error } = await authClient.organization.setActiveTeam({    teamId: "team-id",});

| Prop | Description | Type |
| --- | --- | --- |
| `teamId?` | 
The team ID of the team to set as the current active team.

 | `string` |

#### [List User Teams](#list-user-teams)

List all teams that the current user is a part of.

GET

/organization/list-user-teams

    const { data, error } = await authClient.organization.listUserTeams();

#### [List Team Members](#list-team-members)

List the members of the given team.

POST

/organization/list-team-members

    const { data, error } = await authClient.organization.listTeamMembers({    query: {        teamId: "team-id",    },});

| Prop | Description | Type |
| --- | --- | --- |
| `teamId?` | 
The team whose members we should return. If this is not provided the members of the current active team get returned.

 | `string` |

#### [Add Team Member](#add-team-member)

Add a member to a team.

POST

/organization/add-team-member

    const { data, error } = await authClient.organization.addTeamMember({    teamId: "team-id", // required    userId: "user-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `teamId` | 
The team the user should be a member of.

 | `string` |
| `userId` | 

The user ID which represents the user to be added as a member.

 | `string` |

#### [Remove Team Member](#remove-team-member)

Remove a member from a team.

POST

/organization/remove-team-member

    const { data, error } = await authClient.organization.removeTeamMember({    teamId: "team-id", // required    userId: "user-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `teamId` | 
The team the user should be removed from.

 | `string` |
| `userId` | 

The user which should be removed from the team.

 | `string` |

### [Team Permissions](#team-permissions)

Teams follow the organization's permission system. To manage teams, users need the following permissions:

*   `team:create` - Create new teams
*   `team:update` - Update team details
*   `team:delete` - Remove teams

By default:

*   Organization owners and admins can manage teams
*   Regular members cannot create, update, or delete teams

### [Team Configuration Options](#team-configuration-options)

The teams feature supports several configuration options:

*   `maximumTeams`: Limit the number of teams per organization
    
        teams: {
          enabled: true,
          maximumTeams: 10 // Fixed number
          // OR
          maximumTeams: async ({ organizationId, session }, ctx) => {
            // Dynamic limit based on organization plan
            const plan = await getPlan(organizationId)
            return plan === 'pro' ? 20 : 5
          },
          maximumMembersPerTeam: 10 // Fixed number
          // OR
          maximumMembersPerTeam: async ({ teamId, session, organizationId }, ctx) => {
            // Dynamic limit based on team plan
            const plan = await getPlan(organizationId, teamId)
            return plan === 'pro' ? 50 : 10
          },
        }
    
*   `allowRemovingAllTeams`: Control whether the last team can be removed
    
        teams: {
          enabled: true,
          allowRemovingAllTeams: false // Prevent removing the last team
        }
    

### [Team Members](#team-members)

When inviting members to an organization, you can specify a team:

    await authClient.organization.inviteMember({
      email: "user@example.com",
      role: "member",
      teamId: "team-id",
    });

The invited member will be added to the specified team upon accepting the invitation.

### [Database Schema](#database-schema)

When teams are enabled, new `team` and `teamMember` tables are added to the database.

Table Name: `team`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each team |
| name | string | \- | The name of the team |
| organizationId | string |  | The ID of the organization |
| createdAt | Date | \- | Timestamp of when the team was created |
| updatedAt | Date |  | Timestamp of when the team was created |

Table Name: `teamMember`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each team member |
| teamId | string |  | Unique identifier for each team |
| userId | string |  | The ID of the user |
| createdAt | Date | \- | Timestamp of when the team member was created |

The organization plugin adds the following tables to the database:

### [Organization](#organization-1)

Table Name: `organization`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each organization |
| name | string | \- | The name of the organization |
| slug | string | \- | The slug of the organization |
| logo | string |  | The logo of the organization |
| metadata | string |  | Additional metadata for the organization |
| createdAt | Date | \- | Timestamp of when the organization was created |

### [Member](#member)

Table Name: `member`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each member |
| userId | string |  | The ID of the user |
| organizationId | string |  | The ID of the organization |
| role | string | \- | The role of the user in the organization |
| createdAt | Date | \- | Timestamp of when the member was added to the organization |

### [Invitation](#invitation)

Table Name: `invitation`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each invitation |
| email | string | \- | The email address of the user |
| inviterId | string |  | The ID of the inviter |
| organizationId | string |  | The ID of the organization |
| role | string | \- | The role of the user in the organization |
| status | string | \- | The status of the invitation |
| createdAt | Date | \- | Timestamp of when the invitation was created |
| expiresAt | Date | \- | Timestamp of when the invitation expires |

If teams are enabled, you need to add the following fields to the invitation table:

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| teamId | string |  | The ID of the team |

### [Session](#session)

Table Name: `session`

You need to add two more fields to the session table to store the active organization ID and the active team ID.

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| activeOrganizationId | string |  | The ID of the active organization |
| activeTeamId | string |  | The ID of the active team |

### [Teams (optional)](#teams-optional)

Table Name: `team`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each team |
| name | string | \- | The name of the team |
| organizationId | string |  | The ID of the organization |
| createdAt | Date | \- | Timestamp of when the team was created |
| updatedAt | Date |  | Timestamp of when the team was created |

Table Name: `teamMember`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each team member |
| teamId | string |  | Unique identifier for each team |
| userId | string |  | The ID of the user |
| createdAt | Date | \- | Timestamp of when the team member was created |

Table Name: `invitation`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| teamId | string |  | The ID of the team |

### [Customizing the Schema](#customizing-the-schema)

To change the schema table name or fields, you can pass `schema` option to the organization plugin.

auth.ts

    const auth = betterAuth({
      plugins: [
        organization({
          schema: {
            organization: {
              modelName: "organizations", //map the organization table to organizations
              fields: {
                name: "title", //map the name field to title
              },
              additionalFields: {
                // Add a new field to the organization table
                myCustomField: {
                  type: "string",
                  input: true,
                  required: false,
                },
              },
            },
          },
        }),
      ],
    });

#### [Additional Fields](#additional-fields-1)

Starting with [Better Auth v1.3](https://github.com/better-auth/better-auth/releases/tag/v1.3.0), you can easily add custom fields to the `organization`, `invitation`, `member`, and `team` tables.

When you add extra fields to a model, the relevant API endpoints will automatically accept and return these new properties. For instance, if you add a custom field to the `organization` table, the `createOrganization` endpoint will include this field in its request and response payloads as needed.

auth.ts

    const auth = betterAuth({
      plugins: [
        organization({
          schema: {
            organization: {
              additionalFields: {
                myCustomField: {
                  type: "string", 
                  input: true, 
                  required: false, 
                }, 
              },
            },
          },
        }),
      ],
    });

For inferring the additional fields, you can use the `inferOrgAdditionalFields` function. This function will infer the additional fields from the auth object type.

auth-client.ts

    import { createAuthClient } from "better-auth/client";
    import {
      inferOrgAdditionalFields,
      organizationClient,
    } from "better-auth/client/plugins";
    import type { auth } from "@/auth"; // import the auth object type only
    
    const client = createAuthClient({
      plugins: [
        organizationClient({
          schema: inferOrgAdditionalFields<typeof auth>(),
        }),
      ],
    });

if you can't import the auth object type, you can use the `inferOrgAdditionalFields` function without the generic. This function will infer the additional fields from the schema object.

auth-client.ts

    const client = createAuthClient({
      plugins: [
        organizationClient({
          schema: inferOrgAdditionalFields({
            organization: {
              additionalFields: {
                newField: {
                  type: "string", 
                }, 
              },
            },
          }),
        }),
      ],
    });
    
    //example usage
    await client.organization.create({
      name: "Test",
      slug: "test",
      newField: "123", //this should be allowed
      //@ts-expect-error - this field is not available
      unavailableField: "123", //this should be not allowed
    });

**allowUserToCreateOrganization**: `boolean` | `((user: User) => Promise<boolean> | boolean)` - A function that determines whether a user can create an organization. By default, it's `true`. You can set it to `false` to restrict users from creating organizations.

**organizationLimit**: `number` | `((user: User) => Promise<boolean> | boolean)` - The maximum number of organizations allowed for a user. By default, it's `5`. You can set it to any number you want or a function that returns a boolean.

**creatorRole**: `admin | owner` - The role of the user who creates the organization. By default, it's `owner`. You can set it to `admin`.

**membershipLimit**: `number` - The maximum number of members allowed in an organization. By default, it's `100`. You can set it to any number you want.

**sendInvitationEmail**: `async (data) => Promise<void>` - A function that sends an invitation email to the user.

**invitationExpiresIn** : `number` - How long the invitation link is valid for in seconds. By default, it's 48 hours (2 days).

**cancelPendingInvitationsOnReInvite**: `boolean` - Whether to cancel pending invitations if the user is already invited to the organization. By default, it's `false`.

**invitationLimit**: `number` | `((user: User) => Promise<boolean> | boolean)` - The maximum number of invitations allowed for a user. By default, it's `100`. You can set it to any number you want or a function that returns a boolean.

**requireEmailVerificationOnInvitation**: `boolean` - Whether to require email verification before accepting or rejecting invitations. By default, it's `false`. When enabled, users must have verified their email address before they can accept or reject organization invitations.</content>
</page>

<page>
  <title>API Key | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/api-key</url>
  <content>The API Key plugin allows you to create and manage API keys for your application. It provides a way to authenticate and authorize API requests by verifying API keys.

*   Create, manage, and verify API keys
*   [Built-in rate limiting](https://www.better-auth.com/docs/plugins/api-key#rate-limiting)
*   [Custom expiration times, remaining count, and refill systems](https://www.better-auth.com/docs/plugins/api-key#remaining-refill-and-expiration)
*   [metadata for API keys](https://www.better-auth.com/docs/plugins/api-key#metadata)
*   Custom prefix
*   [Sessions from API keys](https://www.better-auth.com/docs/plugins/api-key#sessions-from-api-keys)
*   [Secondary storage support](https://www.better-auth.com/docs/plugins/api-key#secondary-storage) for high-performance API key lookups

### [Add Plugin to the server](#add-plugin-to-the-server)

auth.ts

    import { betterAuth } from "better-auth"
    import { apiKey } from "better-auth/plugins"
    
    export const auth = betterAuth({
        plugins: [ 
            apiKey() 
        ] 
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary fields and tables to the database.

See the [Schema](#schema) section to add the fields manually.

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { apiKeyClient } from "better-auth/client/plugins"
    
    export const authClient = createAuthClient({
        plugins: [ 
            apiKeyClient() 
        ] 
    })

You can view the list of API Key plugin options [here](https://www.better-auth.com/docs/plugins/api-key#api-key-plugin-options).

### [Create an API key](#create-an-api-key)

Notes

You can adjust more specific API key configurations by using the server method instead.

    const { data, error } = await authClient.apiKey.create({    name: 'project-api-key',    expiresIn: 60 * 60 * 24 * 7,    prefix: 'project-api-key',    metadata: { someKey: 'someValue' },});

| Prop | Description | Type |
| --- | --- | --- |
| `name?` | 
Name of the Api Key.

 | `string` |
| `expiresIn?` | 

Expiration time of the Api Key in seconds.

 | `number` |
| `prefix?` | 

Prefix of the Api Key.

 | `string` |
| `metadata?` | 

Metadata of the Api Key.

 | `any | null` |

API keys are assigned to a user.

#### [Result](#result)

It'll return the `ApiKey` object which includes the `key` value for you to use. Otherwise if it throws, it will throw an `APIError`.

* * *

### [Verify an API key](#verify-an-api-key)

    const permissions = { // Permissions to check are optional.  projects: ["read", "read-write"],}const data = await auth.api.verifyApiKey({    body: {        key: "your_api_key_here", // required        permissions,    },});

| Prop | Description | Type |
| --- | --- | --- |
| `key` | 
The key to verify.

 | `string` |
| `permissions?` | 

The permissions to verify. Optional.

 | `Record<string, string[]>` |

#### [Result](#result-1)

    type Result = {
      valid: boolean;
      error: { message: string; code: string } | null;
      key: Omit<ApiKey, "key"> | null;
    };

* * *

### [Get an API key](#get-an-api-key)

    const { data, error } = await authClient.apiKey.get({    query: {        id: "some-api-key-id", // required    },});

| Prop | Description | Type |
| --- | --- | --- |
| `id` | 
The id of the Api Key.

 | `string` |

#### [Result](#result-2)

You'll receive everything about the API key details, except for the `key` value itself. If it fails, it will throw an `APIError`.

    type Result = Omit<ApiKey, "key">;

* * *

### [Update an API key](#update-an-api-key)

    const { data, error } = await authClient.apiKey.update({    keyId: "some-api-key-id", // required    name: "some-api-key-name",});

| Prop | Description | Type |
| --- | --- | --- |
| `keyId` | 
The id of the Api Key to update.

 | `string` |
| `name?` | 

The name of the key.

 | `string` |

#### [Result](#result-3)

If fails, throws `APIError`. Otherwise, you'll receive the API Key details, except for the `key` value itself.

* * *

### [Delete an API Key](#delete-an-api-key)

Notes

This endpoint is attempting to delete the API key from the perspective of the user. It will check if the user's ID matches the key owner to be able to delete it. If you want to delete a key without these checks, we recommend you use an ORM to directly mutate your DB instead.

    const { data, error } = await authClient.apiKey.delete({    keyId: "some-api-key-id", // required});

| Prop | Description | Type |
| --- | --- | --- |
| `keyId` | 
The id of the Api Key to delete.

 | `string` |

#### [Result](#result-4)

If fails, throws `APIError`. Otherwise, you'll receive:

    type Result = {
      success: boolean;
    };

* * *

### [List API keys](#list-api-keys)

    const { data, error } = await authClient.apiKey.list();

#### [Result](#result-5)

If fails, throws `APIError`. Otherwise, you'll receive:

    type Result = ApiKey[];

* * *

### [Delete all expired API keys](#delete-all-expired-api-keys)

This function will delete all API keys that have an expired expiration date.

POST

/api-key/delete-all-expired-api-keys

    const data = await auth.api.deleteAllExpiredApiKeys();

We automatically delete expired API keys every time any apiKey plugin endpoints were called, however they are rate-limited to a 10 second cool down each call to prevent multiple calls to the database.

* * *

Any time an endpoint in Better Auth is called that has a valid API key in the headers, you can automatically create a mock session to represent the user by enabling `sessionForAPIKeys` option.

This is generally not recommended, as it can lead to security issues if not used carefully. A leaked api key can be used to impersonate a user.

**Rate Limiting Note**: When `enableSessionForAPIKeys` is enabled, the API key is validated once per request, and rate limiting is applied accordingly. If you manually verify an API key and then fetch a session separately, both operations will increment the rate limit counter. Using `enableSessionForAPIKeys` avoids this double increment.

    export const auth = betterAuth({
      plugins: [
        apiKey({
          enableSessionForAPIKeys: true,
        }),
      ],
    });

The default header key is `x-api-key`, but this can be changed by setting the `apiKeyHeaders` option in the plugin options.

    export const auth = betterAuth({
      plugins: [
        apiKey({
          apiKeyHeaders: ["x-api-key", "xyz-api-key"], // or you can pass just a string, eg: "x-api-key"
        }),
      ],
    });

Or optionally, you can pass an `apiKeyGetter` function to the plugin options, which will be called with the `GenericEndpointContext`, and from there, you should return the API key, or `null` if the request is invalid.

    export const auth = betterAuth({
      plugins: [
        apiKey({
          apiKeyGetter: (ctx) => {
            const has = ctx.request.headers.has("x-api-key");
            if (!has) return null;
            return ctx.request.headers.get("x-api-key");
          },
        }),
      ],
    });

The API Key plugin supports multiple storage modes for flexible API key management, allowing you to choose the best strategy for your use case.

### [Storage Mode Options](#storage-mode-options)

#### [`"database"` (Default)](#database-default)

Store API keys only in the database adapter. This is the default mode and requires no additional configuration.

    export const auth = betterAuth({
      plugins: [
        apiKey({
          storage: "database", // Default, can be omitted
        }),
      ],
    });

#### [`"secondary-storage"`](#secondary-storage)

Store API keys only in secondary storage (e.g., Redis). No fallback to database. Best for high-performance scenarios where all keys are migrated to secondary storage.

    import { createClient } from "redis";
    import { betterAuth } from "better-auth";
    import { apiKey } from "better-auth/plugins";
    
    const redis = createClient();
    await redis.connect();
    
    export const auth = betterAuth({
      secondaryStorage: {
        get: async (key) => await redis.get(key),
        set: async (key, value, ttl) => {
          if (ttl) await redis.set(key, value, { EX: ttl });
          else await redis.set(key, value);
        },
        delete: async (key) => await redis.del(key),
      },
      plugins: [
        apiKey({
          storage: "secondary-storage",
        }),
      ],
    });

#### [Secondary Storage with Fallback](#secondary-storage-with-fallback)

Check secondary storage first, then fallback to database if not found.

**Read behavior:**

*   Checks secondary storage first
*   If not found, queries the database
*   **Automatically populates secondary storage** when falling back to database (cache warming)
*   Ensures frequently accessed keys stay in cache over time

**Write behavior:**

*   Writes to **both** database and secondary storage
*   Ensures consistency between both storage layers

    export const auth = betterAuth({
      secondaryStorage: {
        get: async (key) => await redis.get(key),
        set: async (key, value, ttl) => {
          if (ttl) await redis.set(key, value, { EX: ttl });
          else await redis.set(key, value);
        },
        delete: async (key) => await redis.del(key),
      },
      plugins: [
        apiKey({
          storage: "secondary-storage",
          fallbackToDatabase: true,
        }),
      ],
    });

### [Custom Storage Methods](#custom-storage-methods)

You can provide custom storage methods specifically for API keys, overriding the global `secondaryStorage` configuration:

    export const auth = betterAuth({
      plugins: [
        apiKey({
          storage: "secondary-storage",
          customStorage: {
            get: async (key) => {
              // Custom get logic for API keys
              return await customStorage.get(key);
            },
            set: async (key, value, ttl) => {
              // Custom set logic for API keys
              await customStorage.set(key, value, ttl);
            },
            delete: async (key) => {
              // Custom delete logic for API keys
              await customStorage.delete(key);
            },
          },
        }),
      ],
    });

Every API key can have its own rate limit settings. The built-in rate-limiting applies whenever an API key is validated, which includes:

*   When verifying an API key via the `/api-key/verify` endpoint
*   When using API keys for session creation (if `enableSessionForAPIKeys` is enabled), rate limiting applies to all endpoints that use the API key

For other endpoints/methods that don't use API keys, you should utilize Better Auth's [built-in rate-limiting](https://www.better-auth.com/docs/concepts/rate-limit).

**Double Rate-Limit Increment**: If you manually verify an API key using `verifyApiKey()` and then fetch a session using `getSession()` with the same API key header, both operations will increment the rate limit counter, resulting in two increments for a single request. To avoid this, either:

*   Use `enableSessionForAPIKeys: true` and let Better Auth handle session creation automatically (recommended)
*   Or verify the API key once and reuse the validated result instead of calling both methods separately

You can refer to the rate-limit default configurations below in the API Key plugin options.

An example default value:

    export const auth = betterAuth({
      plugins: [
        apiKey({
          rateLimit: {
            enabled: true,
            timeWindow: 1000 * 60 * 60 * 24, // 1 day
            maxRequests: 10, // 10 requests per day
          },
        }),
      ],
    });

For each API key, you can customize the rate-limit options on create.

You can only customize the rate-limit options on the server auth instance.

    const apiKey = await auth.api.createApiKey({
      body: {
        rateLimitEnabled: true,
        rateLimitTimeWindow: 1000 * 60 * 60 * 24, // 1 day
        rateLimitMax: 10, // 10 requests per day
      },
      headers: user_headers,
    });

### [How does it work?](#how-does-it-work)

The rate limiting system uses a sliding window approach:

1.  **First Request**: When an API key is used for the first time (no previous `lastRequest`), the request is allowed and `requestCount` is set to 1.
    
2.  **Within Window**: For subsequent requests within the `timeWindow`, the `requestCount` is incremented. If `requestCount` reaches `rateLimitMax`, the request is rejected with a `RATE_LIMITED` error code.
    
3.  **Window Reset**: If the time since the last request exceeds the `timeWindow`, the window resets: `requestCount` is reset to 1 and `lastRequest` is updated to the current time.
    
4.  **Rate Limit Exceeded**: When a request is rejected due to rate limiting, the error response includes a `tryAgainIn` value (in milliseconds) indicating how long to wait before the window resets.
    

**Disabling Rate Limiting**:

*   **Globally**: Set `rateLimit.enabled: false` in plugin options
*   **Per Key**: Set `rateLimitEnabled: false` when creating or updating an API key
*   **Null Values**: If `rateLimitTimeWindow` or `rateLimitMax` is `null`, rate limiting is effectively disabled for that key

When rate limiting is disabled (globally or per-key), requests are still allowed but `lastRequest` is updated for tracking purposes.

[Remaining, refill, and expiration](#remaining-refill-and-expiration)
---------------------------------------------------------------------

The remaining count is the number of requests left before the API key is disabled. The refill interval is the interval in milliseconds where the `remaining` count is refilled when the interval has passed since the last refill (or since creation if no refill has occurred yet). The expiration time is the expiration date of the API key.

### [How does it work?](#how-does-it-work-1)

#### [Remaining:](#remaining)

Whenever an API key is used, the `remaining` count is updated. If the `remaining` count is `null`, then there is no cap to key usage. Otherwise, the `remaining` count is decremented by 1. If the `remaining` count is 0, then the API key is disabled & removed.

#### [refillInterval & refillAmount:](#refillinterval--refillamount)

Whenever an API key is created, the `refillInterval` and `refillAmount` are set to `null` by default. This means that the API key will not be refilled automatically. However, if both `refillInterval` & `refillAmount` are set, then whenever the API key is used:

*   The system checks if the time since the last refill (or since creation if no refill has occurred) exceeds the `refillInterval`
*   If the interval has passed, the `remaining` count is reset to `refillAmount` (not incremented)
*   The `lastRefillAt` timestamp is updated to the current time

#### [Expiration:](#expiration)

Whenever an API key is created, the `expiresAt` is set to `null`. This means that the API key will never expire. However, if the `expiresIn` is set, then the API key will expire after the `expiresIn` time.

You can customize the key generation and verification process straight from the plugin options.

Here's an example:

    export const auth = betterAuth({
      plugins: [
        apiKey({
          customKeyGenerator: (options: {
            length: number;
            prefix: string | undefined;
          }) => {
            const apiKey = mySuperSecretApiKeyGenerator(
              options.length,
              options.prefix
            );
            return apiKey;
          },
          customAPIKeyValidator: async ({ ctx, key }) => {
            const res = await keyService.verify(key)
            return res.valid
          },
        }),
      ],
    });

If you're **not** using the `length` property provided by `customKeyGenerator`, you **must** set the `defaultKeyLength` property to how long generated keys will be.

    export const auth = betterAuth({
      plugins: [
        apiKey({
          customKeyGenerator: () => {
            return crypto.randomUUID();
          },
          defaultKeyLength: 36, // Or whatever the length is
        }),
      ],
    });

If an API key is validated from your `customAPIKeyValidator`, we still must match that against the database's key. However, by providing this custom function, you can improve the performance of the API key verification process, as all failed keys can be invalidated without having to query your database.

We allow you to store metadata alongside your API keys. This is useful for storing information about the key, such as a subscription plan for example.

To store metadata, make sure you haven't disabled the metadata feature in the plugin options.

    export const auth = betterAuth({
      plugins: [
        apiKey({
          enableMetadata: true,
        }),
      ],
    });

Then, you can store metadata in the `metadata` field of the API key object.

    const apiKey = await auth.api.createApiKey({
      body: {
        metadata: {
          plan: "premium",
        },
      },
    });

You can then retrieve the metadata from the API key object.

    const apiKey = await auth.api.getApiKey({
      body: {
        keyId: "your_api_key_id_here",
      },
    });
    
    console.log(apiKey.metadata.plan); // "premium"

`apiKeyHeaders` `string | string[];`

The header name to check for API key. Default is `x-api-key`.

`customAPIKeyGetter` `(ctx: GenericEndpointContext) => string | null`

A custom function to get the API key from the context.

`customAPIKeyValidator` `(options: { ctx: GenericEndpointContext; key: string; }) => boolean | Promise<boolean>`

A custom function to validate the API key.

`customKeyGenerator` `(options: { length: number; prefix: string | undefined; }) => string | Promise<string>`

A custom function to generate the API key.

`startingCharactersConfig` `{ shouldStore?: boolean; charactersLength?: number; }`

Customize the starting characters configuration.

`defaultKeyLength` `number`

The length of the API key. Longer is better. Default is 64. (Doesn't include the prefix length)

`defaultPrefix` `string`

The prefix of the API key.

Note: We recommend you append an underscore to the prefix to make the prefix more identifiable. (eg `hello_`)

`maximumPrefixLength` `number`

The maximum length of the prefix.

`minimumPrefixLength` `number`

The minimum length of the prefix.

`requireName` `boolean`

Whether to require a name for the API key. Default is `false`.

`maximumNameLength` `number`

The maximum length of the name.

`minimumNameLength` `number`

The minimum length of the name.

`enableMetadata` `boolean`

Whether to enable metadata for an API key.

`keyExpiration` `{ defaultExpiresIn?: number | null; disableCustomExpiresTime?: boolean; minExpiresIn?: number; maxExpiresIn?: number; }`

Customize the key expiration.

`rateLimit` `{ enabled?: boolean; timeWindow?: number; maxRequests?: number; }`

Customize the rate-limiting.

`schema` `InferOptionSchema<ReturnType<typeof apiKeySchema>>`

Custom schema for the API key plugin.

`enableSessionForAPIKeys` `boolean`

An API Key can represent a valid session, so we can mock a session for the user if we find a valid API key in the request headers. Default is `false`.

`storage` `"database" | "secondary-storage"`

Storage backend for API keys. Default is `"database"`.

*   `"database"`: Store API keys in the database adapter (default)
*   `"secondary-storage"`: Store API keys in the configured secondary storage (e.g., Redis)

`fallbackToDatabase` `boolean`

When `storage` is `"secondary-storage"`, enable fallback to database if key is not found in secondary storage. Default is `false`.

When `storage` is set to `"secondary-storage"`, you must configure `secondaryStorage` in your Better Auth options. API keys will be stored using key-value patterns:

*   `api-key:${hashedKey}` - Primary lookup by hashed key
*   `api-key:by-id:${id}` - Lookup by ID
*   `api-key:by-user:${userId}` - User's API key list

If an API key has an expiration date (`expiresAt`), a TTL will be automatically set in secondary storage to ensure automatic cleanup.

    export const auth = betterAuth({
      secondaryStorage: {
        get: async (key) => {
          return await redis.get(key);
        },
        set: async (key, value, ttl) => {
          if (ttl) await redis.set(key, value, { EX: ttl });
          else await redis.set(key, value);
        },
        delete: async (key) => {
          await redis.del(key);
        },
      },
      plugins: [
        apiKey({
          storage: "secondary-storage",
        }),
      ],
    });

`customStorage` `{ get: (key: string) => Promise<unknown> | unknown; set: (key: string, value: string, ttl?: number) => Promise<void | null | unknown> | void; delete: (key: string) => Promise<void | null | string> | void; }`

Custom storage methods for API keys. If provided, these methods will be used instead of `ctx.context.secondaryStorage`. Custom methods take precedence over global secondary storage.

Useful when you want to use a different storage backend specifically for API keys, or when you need custom logic for storage operations.

    export const auth = betterAuth({
      plugins: [
        apiKey({
          storage: "secondary-storage",
          customStorage: {
            get: async (key) => await customStorage.get(key),
            set: async (key, value, ttl) => await customStorage.set(key, value, ttl),
            delete: async (key) => await customStorage.delete(key), 
          },
        }),
      ],
    });

`permissions` `{ defaultPermissions?: Statements | ((userId: string, ctx: GenericEndpointContext) => Statements | Promise<Statements>) }`

Permissions for the API key.

Read more about permissions [here](https://www.better-auth.com/docs/plugins/api-key#permissions).

`disableKeyHashing` `boolean`

Disable hashing of the API key.

⚠️ Security Warning: It's strongly recommended to not disable hashing. Storing API keys in plaintext makes them vulnerable to database breaches, potentially exposing all your users' API keys.

* * *

API keys can have permissions associated with them, allowing you to control access at a granular level. Permissions are structured as a record of resource types to arrays of allowed actions.

### [Setting Default Permissions](#setting-default-permissions)

You can configure default permissions that will be applied to all newly created API keys:

    export const auth = betterAuth({
      plugins: [
        apiKey({
          permissions: {
            defaultPermissions: {
              files: ["read"],
              users: ["read"],
            },
          },
        }),
      ],
    });

You can also provide a function that returns permissions dynamically:

    export const auth = betterAuth({
      plugins: [
        apiKey({
          permissions: {
            defaultPermissions: async (userId, ctx) => {
              // Fetch user role or other data to determine permissions
              return {
                files: ["read"],
                users: ["read"],
              };
            },
          },
        }),
      ],
    });

### [Creating API Keys with Permissions](#creating-api-keys-with-permissions)

When creating an API key, you can specify custom permissions:

    const apiKey = await auth.api.createApiKey({
      body: {
        name: "My API Key",
        permissions: {
          files: ["read", "write"],
          users: ["read"],
        },
        userId: "userId",
      },
    });

### [Verifying API Keys with Required Permissions](#verifying-api-keys-with-required-permissions)

When verifying an API key, you can check if it has the required permissions:

    const result = await auth.api.verifyApiKey({
      body: {
        key: "your_api_key_here",
        permissions: {
          files: ["read"],
        },
      },
    });
    
    if (result.valid) {
      // API key is valid and has the required permissions
    } else {
      // API key is invalid or doesn't have the required permissions
    }

### [Updating API Key Permissions](#updating-api-key-permissions)

You can update the permissions of an existing API key:

    const apiKey = await auth.api.updateApiKey({
      body: {
        keyId: existingApiKeyId,
        permissions: {
          files: ["read", "write", "delete"],
          users: ["read", "write"],
        },
      },
      headers: user_headers,
    });

### [Permissions Structure](#permissions-structure)

Permissions follow a resource-based structure:

    type Permissions = {
      [resourceType: string]: string[];
    };
    
    // Example:
    const permissions = {
      files: ["read", "write", "delete"],
      users: ["read"],
      projects: ["read", "write"],
    };

When verifying an API key, all required permissions must be present in the API key's permissions for validation to succeed.

Table: `apiKey`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | The ID of the API key. |
| name | string |  | The name of the API key. |
| start | string |  | The starting characters of the API key. Useful for showing the first few characters of the API key in the UI for the users to easily identify. |
| prefix | string |  | The API Key prefix. Stored as plain text. |
| key | string | \- | The hashed API key itself. |
| userId | string |  | The ID of the user associated with the API key. |
| refillInterval | number |  | The interval to refill the key in milliseconds. |
| refillAmount | number |  | The amount to refill the remaining count of the key. |
| lastRefillAt | Date |  | The date and time when the key was last refilled. |
| enabled | boolean | \- | Whether the API key is enabled. |
| rateLimitEnabled | boolean | \- | Whether the API key has rate limiting enabled. |
| rateLimitTimeWindow | number |  | The time window in milliseconds for the rate limit. |
| rateLimitMax | number |  | The maximum number of requests allowed within the \`rateLimitTimeWindow\`. |
| requestCount | number | \- | The number of requests made within the rate limit time window. |
| remaining | number |  | The number of requests remaining. |
| lastRequest | Date |  | The date and time of the last request made to the key. |
| expiresAt | Date |  | The date and time when the key will expire. |
| createdAt | Date | \- | The date and time the API key was created. |
| updatedAt | Date | \- | The date and time the API key was updated. |
| permissions | string |  | The permissions of the key. |
| metadata | Object |  | Any additional metadata you want to store with the key. |</content>
</page>

<page>
  <title>Dodo Payments | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/dodopayments</url>
  <content>[Dodo Payments](https://dodopayments.com/) is a global Merchant-of-Record platform that lets AI, SaaS and digital businesses sell in 150+ countries without touching tax, fraud, or compliance. A single, developer-friendly API powers checkout, billing, and payouts so you can launch worldwide in minutes.

[

### Get support on Dodo Payments' Discord

This plugin is maintained by the Dodo Payments team.  
Have questions? Our team is available on Discord to assist you anytime.

](https://discord.gg/bYqAp4ayYh)

*   Automatic customer creation on sign-up
*   Type-safe checkout flows with product slug mapping
*   Self-service customer portal
*   Real-time webhook event processing with signature verification

[

### Get started with Dodo Payments

You need a Dodo Payments account and API keys to use this integration.

](https://app.dodopayments.com/)

Run the following command in your project root:

    npm install @dodopayments/better-auth dodopayments better-auth zod

Add these to your `.env` file:

    DODO_PAYMENTS_API_KEY=your_api_key_here
    DODO_PAYMENTS_WEBHOOK_SECRET=your_webhook_secret_here

Create or update `src/lib/auth.ts`:

    import { betterAuth } from "better-auth";
    import {
      dodopayments,
      checkout,
      portal,
      webhooks,
    } from "@dodopayments/better-auth";
    import DodoPayments from "dodopayments";
    
    export const dodoPayments = new DodoPayments({
      bearerToken: process.env.DODO_PAYMENTS_API_KEY!,
      environment: "test_mode"
    });
    
    export const auth = betterAuth({
      plugins: [
        dodopayments({
          client: dodoPayments,
          createCustomerOnSignUp: true,
          use: [
            checkout({
              products: [
                {
                  productId: "pdt_xxxxxxxxxxxxxxxxxxxxx",
                  slug: "premium-plan",
                },
              ],
              successUrl: "/dashboard/success",
              authenticatedUsersOnly: true,
            }),
            portal(),
            webhooks({
              webhookKey: process.env.DODO_PAYMENTS_WEBHOOK_SECRET!,
              onPayload: async (payload) => {
                console.log("Received webhook:", payload.event_type);
              },
            }),
          ],
        }),
      ],
    });

Set `environment` to `live_mode` for production.

Create or update `src/lib/auth-client.ts`:

    import { dodopaymentsClient } from "@dodopayments/better-auth";
    
    export const authClient = createAuthClient({
      baseURL: process.env.BETTER_AUTH_URL || "http://localhost:3000",
      plugins: [dodopaymentsClient()],
    });

### [Creating a Checkout Session](#creating-a-checkout-session)

    const { data: checkout, error } = await authClient.dodopayments.checkout({
      slug: "premium-plan",
      customer: {
        email: "customer@example.com",
        name: "John Doe",
      },
      billing: {
        city: "San Francisco",
        country: "US",
        state: "CA",
        street: "123 Market St",
        zipcode: "94103",
      },
      referenceId: "order_123",
    });
    
    if (checkout) {
      window.location.href = checkout.url;
    }

### [Accessing the Customer Portal](#accessing-the-customer-portal)

    const { data: customerPortal, error } = await authClient.dodopayments.customer.portal();
    if (customerPortal && customerPortal.redirect) {
      window.location.href = customerPortal.url;
    }

### [Listing Customer Data](#listing-customer-data)

    // Get subscriptions
    const { data: subscriptions, error } =
      await authClient.dodopayments.customer.subscriptions.list({
        query: {
          limit: 10,
          page: 1,
          active: true,
        },
      });
    
    // Get payment history
    const { data: payments, error } = await authClient.dodopayments.customer.payments.list({
      query: {
        limit: 10,
        page: 1,
        status: "succeeded",
      },
    });

### [Webhooks](#webhooks)

The webhooks plugin processes real-time payment events from Dodo Payments with secure signature verification. The default endpoint is `/api/auth/dodopayments/webhooks`.

Generate a webhook secret for your endpoint URL (e.g., `https://your-domain.com/api/auth/dodopayments/webhooks`) in the Dodo Payments Dashboard and set it in your .env file:

    DODO_PAYMENTS_WEBHOOK_SECRET=your_webhook_secret_here

Example handler:

    webhooks({
      webhookKey: process.env.DODO_PAYMENTS_WEBHOOK_SECRET!,
      onPayload: async (payload) => {
        console.log("Received webhook:", payload.event_type);
      },
    });

### [Plugin Options](#plugin-options)

*   **client** (required): DodoPayments client instance
*   **createCustomerOnSignUp** (optional): Auto-create customers on user signup
*   **use** (required): Array of plugins to enable (checkout, portal, webhooks)

### [Checkout Plugin Options](#checkout-plugin-options)

*   **products**: Array of products or async function returning products
*   **successUrl**: URL to redirect after successful payment
*   **authenticatedUsersOnly**: Require user authentication (default: false)

If you encounter any issues, please refer to the [Dodo Payments documentation](https://docs.dodopayments.com/) for troubleshooting steps.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/dodopayments.mdx)</content>
</page>

<page>
  <title>MCP | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/mcp</url>
  <content>`OAuth` `MCP`

The **MCP** plugin lets your app act as an OAuth provider for MCP clients. It handles authentication and makes it easy to issue and manage access tokens for MCP applications.

### [Add the Plugin](#add-the-plugin)

Add the MCP plugin to your auth configuration and specify the login page path.

auth.ts

    import { betterAuth } from "better-auth";
    import { mcp } from "better-auth/plugins";
    
    export const auth = betterAuth({
        plugins: [
            mcp({
                loginPage: "/sign-in" // path to your login page
            })
        ]
    });

This doesn't have a client plugin, so you don't need to make any changes to your authClient.

### [Generate Schema](#generate-schema)

Run the migration or generate the schema to add the necessary fields and tables to the database.

The MCP plugin uses the same schema as the OIDC Provider plugin. See the [OIDC Provider Schema](https://www.better-auth.com/docs/plugins/oidc-provider#schema) section for details.

### [OAuth Discovery Metadata](#oauth-discovery-metadata)

Better Auth already handles the `/api/auth/.well-known/oauth-authorization-server` route automatically but some client may fail to parse the `WWW-Authenticate` header and default to `/.well-known/oauth-authorization-server` (this can happen, for example, if your CORS configuration doesn't expose the `WWW-Authenticate`). For this reason it's better to add a route to expose OAuth metadata for MCP clients:

.well-known/oauth-authorization-server/route.ts

    import { oAuthDiscoveryMetadata } from "better-auth/plugins";
    import { auth } from "../../../lib/auth";
    
    export const GET = oAuthDiscoveryMetadata(auth);

### [OAuth Protected Resource Metadata](#oauth-protected-resource-metadata)

Better Auth already handles the `/api/auth/.well-known/oauth-protected-resource` route automatically but some client may fail to parse the `WWW-Authenticate` header and default to `/.well-known/oauth-protected-resource` (this can happen, for example, if your CORS configuration doesn't expose the `WWW-Authenticate`). For this reason it's better to add a route to expose OAuth metadata for MCP clients:

/.well-known/oauth-protected-resource/route.ts

    import { oAuthProtectedResourceMetadata } from "better-auth/plugins";
    import { auth } from "@/lib/auth";
    
    export const GET = oAuthProtectedResourceMetadata(auth);

### [MCP Session Handling](#mcp-session-handling)

You can use the helper function `withMcpAuth` to get the session and handle unauthenticated calls automatically.

api/\[transport\]/route.ts

    import { auth } from "@/lib/auth";
    import { createMcpHandler } from "@vercel/mcp-adapter";
    import { withMcpAuth } from "better-auth/plugins";
    import { z } from "zod";
    
    const handler = withMcpAuth(auth, (req, session) => {
        // session contains the access token record with scopes and user ID
        return createMcpHandler(
            (server) => {
                server.tool(
                    "echo",
                    "Echo a message",
                    { message: z.string() },
                    async ({ message }) => {
                        return {
                            content: [{ type: "text", text: `Tool echo: ${message}` }],
                        };
                    },
                );
            },
            {
                capabilities: {
                    tools: {
                        echo: {
                            description: "Echo a message",
                        },
                    },
                },
            },
            {
                redisUrl: process.env.REDIS_URL,
                basePath: "/api",
                verboseLogs: true,
                maxDuration: 60,
            },
        )(req);
    });
    
    export { handler as GET, handler as POST, handler as DELETE };

You can also use `auth.api.getMcpSession` to get the session using the access token sent from the MCP client:

api/\[transport\]/route.ts

    import { auth } from "@/lib/auth";
    import { createMcpHandler } from "@vercel/mcp-adapter";
    import { z } from "zod";
    
    const handler = async (req: Request) => {
         // session contains the access token record with scopes and user ID
        const session = await auth.api.getMcpSession({
            headers: req.headers
        })
        if(!session){
            //this is important and you must return 401
            return new Response(null, {
                status: 401
            })
        }
        return createMcpHandler(
            (server) => {
                server.tool(
                    "echo",
                    "Echo a message",
                    { message: z.string() },
                    async ({ message }) => {
                        return {
                            content: [{ type: "text", text: `Tool echo: ${message}` }],
                        };
                    },
                );
            },
            {
                capabilities: {
                    tools: {
                        echo: {
                            description: "Echo a message",
                        },
                    },
                },
            },
            {
                redisUrl: process.env.REDIS_URL,
                basePath: "/api",
                verboseLogs: true,
                maxDuration: 60,
            },
        )(req);
    }
    
    export { handler as GET, handler as POST, handler as DELETE };

The MCP plugin accepts the following configuration options:

### [OIDC Configuration](#oidc-configuration)

The plugin supports additional OIDC configuration options through the `oidcConfig` parameter:

The MCP plugin uses the same schema as the OIDC Provider plugin. See the [OIDC Provider Schema](https://www.better-auth.com/docs/plugins/oidc-provider#schema) section for details.

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/mcp.mdx)</content>
</page>

<page>
  <title>Migrating from Supabase Auth to Better Auth | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/supabase-migration-guide</url>
  <content>    import { generateId } from 'better-auth';
    import { DBFieldAttribute } from 'better-auth/db';
    import { Pool } from 'pg';
    import { auth } from './auth'; // <- Your Better Auth Instance
    
    // ============================================================================
    // CONFIGURATION
    // ============================================================================
    
    const CONFIG = {
      /**
       * Number of users to process in each batch
       * Higher values = faster migration but more memory usage
       * Recommended: 5000-10000 for most cases
       */
      batchSize: 5000,
      /**
       * Resume from a specific user ID (cursor-based pagination)
       * Useful for resuming interrupted migrations
       * Set to null to start from the beginning
       */
      resumeFromId: null as string | null,
      /**
       * Temporary email domain for phone-only users
       * Phone-only users need an email for Better Auth
       * Format: {phone_number}@{tempEmailDomain}
       */
      tempEmailDomain: 'temp.better-auth.com',
    };
    
    // ============================================================================
    // TYPE DEFINITIONS
    // ============================================================================
    
    type MigrationStatus = 'idle' | 'running' | 'paused' | 'completed' | 'failed';
    
    type MigrationState = {
      status: MigrationStatus;
      totalUsers: number;
      processedUsers: number;
      successCount: number;
      failureCount: number;
      skipCount: number;
      currentBatch: number;
      totalBatches: number;
      startedAt: Date | null;
      completedAt: Date | null;
      lastProcessedId: string | null;
      errors: Array<{ userId: string; error: string }>;
    };
    
    type UserInsertData = {
      id: string;
      email: string | null;
      name: string;
      emailVerified: boolean;
      createdAt: string | null;
      updatedAt: string | null;
      image?: string;
      [key: string]: any;
    };
    
    type AccountInsertData = {
      id: string;
      userId: string;
      providerId: string;
      accountId: string;
      password: string | null;
      createdAt: string | null;
      updatedAt: string | null;
    };
    
    type SupabaseIdentityFromDB = {
      id: string;
      provider_id: string;
      user_id: string;
      identity_data: Record<string, any>;
      provider: string;
      last_sign_in_at: string | null;
      created_at: string | null;
      updated_at: string | null;
      email: string | null;
    };
    
    type SupabaseUserFromDB = {
      instance_id: string | null;
      id: string;
      aud: string | null;
      role: string | null;
      email: string | null;
      encrypted_password: string | null;
      email_confirmed_at: string | null;
      invited_at: string | null;
      confirmation_token: string | null;
      confirmation_sent_at: string | null;
      recovery_token: string | null;
      recovery_sent_at: string | null;
      email_change_token_new: string | null;
      email_change: string | null;
      email_change_sent_at: string | null;
      last_sign_in_at: string | null;
      raw_app_meta_data: Record<string, any> | null;
      raw_user_meta_data: Record<string, any> | null;
      is_super_admin: boolean | null;
      created_at: string | null;
      updated_at: string | null;
      phone: string | null;
      phone_confirmed_at: string | null;
      phone_change: string | null;
      phone_change_token: string | null;
      phone_change_sent_at: string | null;
      confirmed_at: string | null;
      email_change_token_current: string | null;
      email_change_confirm_status: number | null;
      banned_until: string | null;
      reauthentication_token: string | null;
      reauthentication_sent_at: string | null;
      is_sso_user: boolean;
      deleted_at: string | null;
      is_anonymous: boolean;
      identities: SupabaseIdentityFromDB[];
    };
    
    // ============================================================================
    // MIGRATION STATE MANAGER
    // ============================================================================
    
    class MigrationStateManager {
      private state: MigrationState = {
        status: 'idle',
        totalUsers: 0,
        processedUsers: 0,
        successCount: 0,
        failureCount: 0,
        skipCount: 0,
        currentBatch: 0,
        totalBatches: 0,
        startedAt: null,
        completedAt: null,
        lastProcessedId: null,
        errors: [],
      };
    
      start(totalUsers: number, batchSize: number) {
        this.state = {
          status: 'running',
          totalUsers,
          processedUsers: 0,
          successCount: 0,
          failureCount: 0,
          skipCount: 0,
          currentBatch: 0,
          totalBatches: Math.ceil(totalUsers / batchSize),
          startedAt: new Date(),
          completedAt: null,
          lastProcessedId: null,
          errors: [],
        };
      }
    
      updateProgress(
        processed: number,
        success: number,
        failure: number,
        skip: number,
        lastId: string | null,
      ) {
        this.state.processedUsers += processed;
        this.state.successCount += success;
        this.state.failureCount += failure;
        this.state.skipCount += skip;
        this.state.currentBatch++;
        if (lastId) {
          this.state.lastProcessedId = lastId;
        }
      }
    
      addError(userId: string, error: string) {
        if (this.state.errors.length < 100) {
          this.state.errors.push({ userId, error });
        }
      }
    
      complete() {
        this.state.status = 'completed';
        this.state.completedAt = new Date();
      }
    
      fail() {
        this.state.status = 'failed';
        this.state.completedAt = new Date();
      }
    
      getState(): MigrationState {
        return { ...this.state };
      }
    
      getProgress(): number {
        if (this.state.totalUsers === 0) return 0;
        return Math.round((this.state.processedUsers / this.state.totalUsers) * 100);
      }
    
      getETA(): string | null {
        if (!this.state.startedAt || this.state.processedUsers === 0) {
          return null;
        }
    
        const elapsed = Date.now() - this.state.startedAt.getTime();
        const avgTimePerUser = elapsed / this.state.processedUsers;
        const remainingUsers = this.state.totalUsers - this.state.processedUsers;
        const remainingMs = avgTimePerUser * remainingUsers;
    
        const seconds = Math.floor(remainingMs / 1000);
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);
    
        if (hours > 0) {
          return `${hours}h ${minutes % 60}m`;
        } else if (minutes > 0) {
          return `${minutes}m ${seconds % 60}s`;
        } else {
          return `${seconds}s`;
        }
      }
    }
    
    // ============================================================================
    // DATABASE CONNECTIONS
    // ============================================================================
    
    const fromDB = new Pool({
      connectionString: process.env.FROM_DATABASE_URL,
    });
    
    const toDB = new Pool({
      connectionString: process.env.TO_DATABASE_URL,
    });
    
    // ============================================================================
    // BETTER AUTH VALIDATION
    // ============================================================================
    
    /**
     * Validates that the imported Better Auth instance meets migration requirements
     */
    async function validateAuthConfig() {
      const ctx = await auth.$context;
      const errors: string[] = [];
    
      // Check emailAndPassword
      if (!ctx.options.emailAndPassword?.enabled) {
        errors.push('emailAndPassword.enabled must be true');
      }
    
      // Check required plugins
      const requiredPlugins = ['admin', 'anonymous', 'phone-number'];
      const plugins = ctx.options.plugins || [];
      const pluginIds = plugins.map((p: any) => p.id);
    
      for (const required of requiredPlugins) {
        if (!pluginIds.includes(required)) {
          errors.push(`Missing required plugin: ${required}`);
        }
      }
    
      // Check required additional fields
      const additionalFields = ctx.options.user?.additionalFields || {};
      const requiredFields: Record<string, DBFieldAttribute> = {
        userMetadata: { type: 'json', required: false, input: false },
        appMetadata: { type: 'json', required: false, input: false },
        invitedAt: { type: 'date', required: false, input: false },
        lastSignInAt: { type: 'date', required: false, input: false },
      };
    
      for (const [fieldName, expectedConfig] of Object.entries(requiredFields)) {
        const fieldConfig = additionalFields[fieldName];
    
        if (!fieldConfig) {
          errors.push(`Missing required user.additionalFields: ${fieldName}`);
        } else {
          // Validate field configuration
          if (fieldConfig.type !== expectedConfig.type) {
            errors.push(
              `user.additionalFields.${fieldName} must have type: '${expectedConfig.type}' (got '${fieldConfig.type}')`,
            );
          }
          if (fieldConfig.required !== expectedConfig.required) {
            errors.push(
              `user.additionalFields.${fieldName} must have required: ${expectedConfig.required}`,
            );
          }
          if (fieldConfig.input !== expectedConfig.input) {
            errors.push(`user.additionalFields.${fieldName} must have input: ${expectedConfig.input}`);
          }
        }
      }
    
      if (errors.length > 0) {
        console.error('\n🟧 Better Auth Configuration Errors:\n');
        errors.forEach((err) => console.error(`   ${err}`));
        console.error('\n🟧 Please update your Better Auth configuration to include:\n');
        console.error('   1. emailAndPassword: { enabled: true }');
        console.error('   2. plugins: [admin(), anonymous(), phoneNumber()]');
        console.error(
          '   3. user.additionalFields: { userMetadata, appMetadata, invitedAt, lastSignInAt }\n',
        );
        process.exit(1);
      }
    
      return ctx;
    }
    
    // ============================================================================
    // MIGRATION LOGIC
    // ============================================================================
    
    const stateManager = new MigrationStateManager();
    
    let ctxCache: {
      hasAnonymousPlugin: boolean;
      hasAdminPlugin: boolean;
      hasPhoneNumberPlugin: boolean;
      supportedProviders: string[];
    } | null = null;
    
    async function processBatch(
      users: SupabaseUserFromDB[],
      ctx: any,
    ): Promise<{
      success: number;
      failure: number;
      skip: number;
      errors: Array<{ userId: string; error: string }>;
    }> {
      const stats = {
        success: 0,
        failure: 0,
        skip: 0,
        errors: [] as Array<{ userId: string; error: string }>,
      };
    
      if (!ctxCache) {
        ctxCache = {
          hasAdminPlugin: ctx.options.plugins?.some((p: any) => p.id === 'admin') || false,
          hasAnonymousPlugin: ctx.options.plugins?.some((p: any) => p.id === 'anonymous') || false,
          hasPhoneNumberPlugin: ctx.options.plugins?.some((p: any) => p.id === 'phone-number') || false,
          supportedProviders: Object.keys(ctx.options.socialProviders || {}),
        };
      }
    
      const { hasAdminPlugin, hasAnonymousPlugin, hasPhoneNumberPlugin, supportedProviders } = ctxCache;
    
      const validUsersData: Array<{ user: SupabaseUserFromDB; userData: UserInsertData }> = [];
    
      for (const user of users) {
        if (!user.email && !user.phone) {
          stats.skip++;
          continue;
        }
        if (!user.email && !hasPhoneNumberPlugin) {
          stats.skip++;
          continue;
        }
        if (user.deleted_at) {
          stats.skip++;
          continue;
        }
        if (user.banned_until && !hasAdminPlugin) {
          stats.skip++;
          continue;
        }
    
        const getTempEmail = (phone: string) =>
          `${phone.replace(/[^0-9]/g, '')}@${CONFIG.tempEmailDomain}`;
    
        const getName = (): string => {
          if (user.raw_user_meta_data?.name) return user.raw_user_meta_data.name;
          if (user.raw_user_meta_data?.full_name) return user.raw_user_meta_data.full_name;
          if (user.raw_user_meta_data?.username) return user.raw_user_meta_data.username;
          if (user.raw_user_meta_data?.user_name) return user.raw_user_meta_data.user_name;
    
          const firstId = user.identities?.[0];
          if (firstId?.identity_data?.name) return firstId.identity_data.name;
          if (firstId?.identity_data?.full_name) return firstId.identity_data.full_name;
          if (firstId?.identity_data?.username) return firstId.identity_data.username;
          if (firstId?.identity_data?.preferred_username)
            return firstId.identity_data.preferred_username;
    
          if (user.email) return user.email.split('@')[0]!;
          if (user.phone) return user.phone;
    
          return 'Unknown';
        };
    
        const getImage = (): string | undefined => {
          if (user.raw_user_meta_data?.avatar_url) return user.raw_user_meta_data.avatar_url;
          if (user.raw_user_meta_data?.picture) return user.raw_user_meta_data.picture;
          const firstId = user.identities?.[0];
          if (firstId?.identity_data?.avatar_url) return firstId.identity_data.avatar_url;
          if (firstId?.identity_data?.picture) return firstId.identity_data.picture;
          return undefined;
        };
    
        const userData: UserInsertData = {
          id: user.id,
          email: user.email || (user.phone ? getTempEmail(user.phone) : null),
          emailVerified: !!user.email_confirmed_at,
          name: getName(),
          image: getImage(),
          createdAt: user.created_at,
          updatedAt: user.updated_at,
        };
    
        if (hasAnonymousPlugin) userData.isAnonymous = user.is_anonymous;
        if (hasPhoneNumberPlugin && user.phone) {
          userData.phoneNumber = user.phone;
          userData.phoneNumberVerified = !!user.phone_confirmed_at;
        }
    
        if (hasAdminPlugin) {
          userData.role = user.is_super_admin ? 'admin' : user.role || 'user';
          if (user.banned_until) {
            const banExpires = new Date(user.banned_until);
            if (banExpires > new Date()) {
              userData.banned = true;
              userData.banExpires = banExpires;
              userData.banReason = 'Migrated from Supabase (banned)';
            } else {
              userData.banned = false;
            }
          } else {
            userData.banned = false;
          }
        }
    
        if (user.raw_user_meta_data && Object.keys(user.raw_user_meta_data).length > 0) {
          userData.userMetadata = user.raw_user_meta_data;
        }
        if (user.raw_app_meta_data && Object.keys(user.raw_app_meta_data).length > 0) {
          userData.appMetadata = user.raw_app_meta_data;
        }
        if (user.invited_at) userData.invitedAt = user.invited_at;
        if (user.last_sign_in_at) userData.lastSignInAt = user.last_sign_in_at;
    
        validUsersData.push({ user, userData });
      }
    
      if (validUsersData.length === 0) {
        return stats;
      }
    
      try {
        await toDB.query('BEGIN');
    
        const allFields = new Set<string>();
        validUsersData.forEach(({ userData }) => {
          Object.keys(userData).forEach((key) => allFields.add(key));
        });
        const fields = Array.from(allFields);
    
        const maxParamsPerQuery = 65000;
        const fieldsPerUser = fields.length;
        const usersPerChunk = Math.floor(maxParamsPerQuery / fieldsPerUser);
    
        for (let i = 0; i < validUsersData.length; i += usersPerChunk) {
          const chunk = validUsersData.slice(i, i + usersPerChunk);
    
          const placeholders: string[] = [];
          const values: any[] = [];
          let paramIndex = 1;
    
          for (const { userData } of chunk) {
            const userPlaceholders = fields.map((field) => {
              values.push(userData[field] ?? null);
              return `$${paramIndex++}`;
            });
            placeholders.push(`(${userPlaceholders.join(', ')})`);
          }
    
          await toDB.query(
            `
            INSERT INTO "user" (${fields.map((f) => `"${f}"`).join(', ')})
            VALUES ${placeholders.join(', ')}
            ON CONFLICT (id) DO NOTHING
          `,
            values,
          );
        }
    
        const accountsData: AccountInsertData[] = [];
    
        for (const { user } of validUsersData) {
          for (const identity of user.identities ?? []) {
            if (identity.provider === 'email') {
              accountsData.push({
                id: generateId(),
                userId: user.id,
                providerId: 'credential',
                accountId: user.id,
                password: user.encrypted_password || null,
                createdAt: user.created_at,
                updatedAt: user.updated_at,
              });
            }
    
            if (supportedProviders.includes(identity.provider)) {
              accountsData.push({
                id: generateId(),
                userId: user.id,
                providerId: identity.provider,
                accountId: identity.identity_data?.sub || identity.provider_id,
                password: null,
                createdAt: identity.created_at ?? user.created_at,
                updatedAt: identity.updated_at ?? user.updated_at,
              });
            }
          }
        }
    
        if (accountsData.length > 0) {
          const maxParamsPerQuery = 65000;
          const fieldsPerAccount = 7;
          const accountsPerChunk = Math.floor(maxParamsPerQuery / fieldsPerAccount);
    
          for (let i = 0; i < accountsData.length; i += accountsPerChunk) {
            const chunk = accountsData.slice(i, i + accountsPerChunk);
    
            const accountPlaceholders: string[] = [];
            const accountValues: any[] = [];
            let paramIndex = 1;
    
            for (const acc of chunk) {
              accountPlaceholders.push(
                `($${paramIndex++}, $${paramIndex++}, $${paramIndex++}, $${paramIndex++}, $${paramIndex++}, $${paramIndex++}, $${paramIndex++})`,
              );
              accountValues.push(
                acc.id,
                acc.userId,
                acc.providerId,
                acc.accountId,
                acc.password,
                acc.createdAt,
                acc.updatedAt,
              );
            }
    
            await toDB.query(
              `
              INSERT INTO "account" ("id", "userId", "providerId", "accountId", "password", "createdAt", "updatedAt")
              VALUES ${accountPlaceholders.join(', ')}
              ON CONFLICT ("id") DO NOTHING
            `,
              accountValues,
            );
          }
        }
    
        await toDB.query('COMMIT');
        stats.success = validUsersData.length;
      } catch (error: any) {
        await toDB.query('ROLLBACK');
        console.error('[TRANSACTION] Batch failed, rolled back:', error.message);
        stats.failure = validUsersData.length;
        if (stats.errors.length < 100) {
          stats.errors.push({ userId: 'bulk', error: error.message });
        }
      }
    
      return stats;
    }
    
    async function migrateFromSupabase() {
      const { batchSize, resumeFromId } = CONFIG;
    
      console.log('[MIGRATION] Starting migration with config:', CONFIG);
    
      // Validate Better Auth configuration
      const ctx = await validateAuthConfig();
    
      try {
        const countResult = await fromDB.query<{ count: string }>(
          `
          SELECT COUNT(*) as count FROM auth.users
          ${resumeFromId ? 'WHERE id > $1' : ''}
        `,
          resumeFromId ? [resumeFromId] : [],
        );
    
        const totalUsers = parseInt(countResult.rows[0]?.count || '0', 10);
    
        console.log(`[MIGRATION] Starting migration for ${totalUsers.toLocaleString()} users`);
        console.log(`[MIGRATION] Batch size: ${batchSize}\n`);
    
        stateManager.start(totalUsers, batchSize);
    
        let lastProcessedId: string | null = resumeFromId;
        let hasMore = true;
        let batchNumber = 0;
    
        while (hasMore) {
          batchNumber++;
          const batchStart = Date.now();
    
          const result: { rows: SupabaseUserFromDB[] } = await fromDB.query<SupabaseUserFromDB>(
            `
            SELECT 
              u.*,
              COALESCE(
                json_agg(
                  i.* ORDER BY i.id
                ) FILTER (WHERE i.id IS NOT NULL),
                '[]'::json
              ) as identities
            FROM auth.users u
            LEFT JOIN auth.identities i ON u.id = i.user_id
            ${lastProcessedId ? 'WHERE u.id > $1' : ''}
            GROUP BY u.id
            ORDER BY u.id ASC
            LIMIT $${lastProcessedId ? '2' : '1'}
          `,
            lastProcessedId ? [lastProcessedId, batchSize] : [batchSize],
          );
    
          const batch: SupabaseUserFromDB[] = result.rows;
          hasMore = batch.length === batchSize;
    
          if (batch.length === 0) break;
    
          console.log(
            `\nBatch ${batchNumber}/${Math.ceil(totalUsers / batchSize)} (${batch.length} users)`,
          );
    
          const stats = await processBatch(batch, ctx);
    
          lastProcessedId = batch[batch.length - 1]!.id;
          stateManager.updateProgress(
            batch.length,
            stats.success,
            stats.failure,
            stats.skip,
            lastProcessedId,
          );
    
          stats.errors.forEach((err) => stateManager.addError(err.userId, err.error));
    
          const batchTime = ((Date.now() - batchStart) / 1000).toFixed(2);
          const usersPerSec = (batch.length / parseFloat(batchTime)).toFixed(0);
    
          const state = stateManager.getState();
          console.log(`Success: ${stats.success} | Skip: ${stats.skip} | Failure: ${stats.failure}`);
          console.log(
            `Progress: ${stateManager.getProgress()}% (${state.processedUsers.toLocaleString()}/${state.totalUsers.toLocaleString()})`,
          );
          console.log(`Speed: ${usersPerSec} users/sec (${batchTime}s for this batch)`);
    
          const eta = stateManager.getETA();
          if (eta) {
            console.log(`ETA: ${eta}`);
          }
        }
    
        stateManager.complete();
        const finalState = stateManager.getState();
    
        console.log('\n🎉 Migration completed');
        console.log(`   - Success: ${finalState.successCount.toLocaleString()}`);
        console.log(`   - Skipped: ${finalState.skipCount.toLocaleString()}`);
        console.log(`   - Failed: ${finalState.failureCount.toLocaleString()}`);
    
        const totalTime =
          finalState.completedAt && finalState.startedAt
            ? ((finalState.completedAt.getTime() - finalState.startedAt.getTime()) / 1000 / 60).toFixed(
                1,
              )
            : '0';
        console.log(`   Total time: ${totalTime} minutes`);
    
        if (finalState.errors.length > 0) {
          console.log(`\nFirst ${Math.min(10, finalState.errors.length)} errors:`);
          finalState.errors.slice(0, 10).forEach((err) => {
            console.log(`- User ${err.userId}: ${err.error}`);
          });
        }
    
        return finalState;
      } catch (error) {
        stateManager.fail();
        console.error('\nMigration failed:', error);
        throw error;
      } finally {
        await fromDB.end();
        await toDB.end();
      }
    }
    
    // ============================================================================
    // MAIN ENTRY POINT
    // ============================================================================
    
    async function main() {
      console.log('🚀 Supabase Auth → Better Auth Migration\n');
    
      if (!process.env.FROM_DATABASE_URL) {
        console.error('Error: FROM_DATABASE_URL environment variable is required');
        process.exit(1);
      }
      if (!process.env.TO_DATABASE_URL) {
        console.error('Error: TO_DATABASE_URL environment variable is required');
        process.exit(1);
      }
    
      try {
        await migrateFromSupabase();
        process.exit(0);
      } catch (error) {
        console.error('\nMigration failed:', error);
        process.exit(1);
      }
    }
    main();</content>
</page>

<page>
  <title>Migrating from Auth.js to Better Auth | Better Auth</title>
  <url>https://www.better-auth.com/docs/guides/next-auth-migration-guide</url>
  <content>In this guide, we'll walk through the steps to migrate a project from [Auth.js](https://authjs.dev/) (formerly [NextAuth.js](https://next-auth.js.org/)) to Better Auth. Since these projects have different design philosophies, the migration requires careful planning and work. If your current setup is working well, there's no urgent need to migrate. We continue to handle security patches and critical issues for Auth.js.

However, if you're starting a new project or facing challenges with your current setup, we strongly recommend using Better Auth. Our roadmap includes features previously exclusive to Auth.js, and we hope this will unite the ecosystem more strongly without causing fragmentation.

Before starting the migration process, set up Better Auth in your project. Follow the [installation guide](https://www.better-auth.com/docs/installation) to get started.

For example, if you use the GitHub OAuth provider, here is a comparison of the `auth.ts` file.

Now Better Auth supports stateless session management without any database. If you were using a Database adapter in Auth.js, you can refer to the [Database models](#database-models) below to check the differences with Better Auth's core schema.

This client instance includes a set of functions for interacting with the Better Auth server instance. For more information, see [here](https://www.better-auth.com/docs/concepts/client).

auth-client.ts

    import { createAuthClient } from "better-auth/react"
    
    export const authClient = createAuthClient()

Rename the `/app/api/auth/[...nextauth]` folder to `/app/api/auth/[...all]` to avoid confusion. Then, update the `route.ts` file as follows:

In this section, we'll look at how to manage sessions in Better Auth compared to Auth.js. For more information, see [here](https://www.better-auth.com/docs/concepts/session-management).

### [Client-side](#client-side)

#### [Sign In](#sign-in)

Here are the differences between Auth.js and Better Auth for signing in users. For example, with the GitHub OAuth provider:

#### [Sign Out](#sign-out)

Here are the differences between Auth.js and Better Auth when signing out users.

#### [Get Session](#get-session)

Here are the differences between Auth.js and Better Auth for getting the current active session.

### [Server-side](#server-side)

#### [Sign In](#sign-in-1)

Here are the differences between Auth.js and Better Auth for signing in users. For example, with the GitHub OAuth provider:

#### [Sign Out](#sign-out-1)

Here are the differences between Auth.js and Better Auth when signing out users.

#### [Get Session](#get-session-1)

Here are the differences between Auth.js and Better Auth for getting the current active session.

> Proxy(Middleware) is not intended for slow data fetching. While Proxy can be helpful for optimistic checks such as permission-based redirects, it should not be used as a full session management or authorization solution. - [Next.js docs](https://nextjs.org/docs/app/getting-started/proxy#use-cases)

Auth.js offers approaches using Proxy (Middleware), but we recommend handling auth checks on each page or route rather than in a Proxy or Layout. Here is a basic example of protecting resources with Better Auth.

Both Auth.js and Better Auth provide stateless (JWT) and database session strategies. If you were using the database session strategy in Auth.js and plan to continue using it in Better Auth, you will also need to migrate your database.

Just like Auth.js has database models, Better Auth also has a core schema. In this section, we'll compare the two and explore the differences between them.

### [User -> User](#user---user)

### [Session -> Session](#session---session)

### [Account -> Account](#account---account)

### [VerificationToken -> Verification](#verificationtoken---verification)

### [Comparison](#comparison)

Table: **User**

*   `name`, `email`, and `emailVerified` are required in Better Auth, while optional in Auth.js
*   `emailVerified` uses a boolean in Better Auth, while Auth.js uses a timestamp
*   Better Auth includes `createdAt` and `updatedAt` timestamps

Table: **Session**

*   Better Auth uses `token` instead of `sessionToken`
*   Better Auth uses `expiresAt` instead of `expires`
*   Better Auth includes `ipAddress` and `userAgent` fields
*   Better Auth includes `createdAt` and `updatedAt` timestamps

Table: **Account**

*   Better Auth uses camelCase naming (e.g. `refreshToken` vs `refresh_token`)
*   Better Auth includes `accountId` to distinguish between the account ID and internal ID
*   Better Auth uses `providerId` instead of `provider`
*   Better Auth includes `accessTokenExpiresAt` and `refreshTokenExpiresAt` for token management
*   Better Auth includes `password` field to support built-in credential authentication
*   Better Auth does not have a `type` field as it's determined by the `providerId`
*   Better Auth removes `token_type` and `session_state` fields
*   Better Auth includes `createdAt` and `updatedAt` timestamps

Table: **VerificationToken -> Verification**

*   Better Auth uses `Verification` table instead of `VerificationToken`
*   Better Auth uses a single `id` primary key instead of composite primary key
*   Better Auth uses `value` instead of `token` to support various verification types
*   Better Auth uses `expiresAt` instead of `expires`
*   Better Auth includes `createdAt` and `updatedAt` timestamps

If you were using Auth.js v4, note that v5 does not introduce any breaking changes to the database schema. Optional fields like `oauth_token_secret` and `oauth_token` can be removed if you are not using them. Rarely used fields like `refresh_token_expires_in` can also be removed.

### [Customization](#customization)

You may have extended the database models or implemented additional logic in Auth.js. Better Auth allows you to customize the core schema in a type-safe way. You can also define custom logic during the lifecycle of database operations. For more details, see [Concepts - Database](https://www.better-auth.com/docs/concepts/database).

Now you're ready to migrate from Auth.js to Better Auth. For a complete implementation with multiple authentication methods, check out the [Next.js Demo App](https://github.com/better-auth/better-auth/tree/canary/demo/nextjs). Better Auth offers greater flexibility and more features, so be sure to explore the [documentation](https://www.better-auth.com/docs) to unlock its full potential.

If you need help with migration, join our [community](https://www.better-auth.com/community) or reach out to [contact@better-auth.com](mailto:contact@better-auth.com).

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/guides/next-auth-migration-guide.mdx)</content>
</page>

<page>
  <title>Autumn Billing | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/autumn</url>
  <content>[Autumn](https://useautumn.com/) is open source infrastructure to run SaaS pricing plans. It sits between your app and Stripe, and acts as the database for your customers' subscription status, usage metering and feature permissions.

[

### Get help on Autumn's Discord

We're online to help you with any questions you have.

](https://discord.gg/STqxY92zuS)

*   One function for all checkout, subscription and payment flows
*   No webhooks required: query Autumn for the data you need
*   Manages your application's free and paid plans
*   Usage tracking for usage billing and periodic limits
*   Custom plans and pricing changes through Autumn's dashboard

### [Setup Autumn Account](#setup-autumn-account)

First, create your pricing plans in Autumn's [dashboard](https://app.useautumn.com/), where you define what each plan and product gets access to and how it should be billed. In this example, we're handling the free and pro plans for an AI chatbot, which comes with a number of `messages` per month.

### [Install Autumn SDK](#install-autumn-sdk)

If you're using a separate client and server setup, make sure to install the plugin in both parts of your project.

### [Add `AUTUMN_SECRET_KEY` to your environment variables](#add-autumn_secret_key-to-your-environment-variables)

You can find it in Autumn's dashboard under "[Developer](https://app.useautumn.com/sandbox/onboarding)".

.env

    AUTUMN_SECRET_KEY=am_sk_xxxxxxxxxx

### [Add the Autumn plugin to your `auth` config](#add-the-autumn-plugin-to-your-auth-config)

Autumn will auto-create your customers when they sign up, and assign them any default plans you created (eg your Free plan). You can choose who becomes a customer: individual users, organizations, both, or something custom like workspaces.

### [Add `<AutumnProvider />`](#add-autumnprovider-)

Client side, wrap your application with the AutumnProvider component, and pass in the `baseUrl` that you define within better-auth's `authClient`.

app/layout.tsx

    import { AutumnProvider } from "autumn-js/react";
    
    export default function RootLayout({
      children,
    }: {
      children: React.ReactNode;
    }) {
      return (
        <html>
          <body>
            {/* or meta.env.BETTER_AUTH_URL for vite */}
            <AutumnProvider betterAuthUrl={process.env.NEXT_PUBLIC_BETTER_AUTH_URL}>
              {children}
            </AutumnProvider>
          </body>
        </html>
      );
    }

### [Handle payments](#handle-payments)

Call `attach` to redirect the customer to a Stripe checkout page when they want to purchase the Pro plan.

If their payment method is already on file, `AttachDialog` will open instead to let the customer confirm their new subscription or purchase, and handle the payment.

    import { useCustomer, AttachDialog } from "autumn-js/react";
    
    export default function PurchaseButton() {
      const { attach } = useCustomer();
    
      return (
        <button
          onClick={async () => {
            await attach({
              productId: "pro",
              dialog: AttachDialog,
            });
          }}
        >
          Upgrade to Pro
        </button>
      );
    }

The AttachDialog component can be used directly from the `autumn-js/react` library (as shown in the example above), or downloaded as a [shadcn/ui component](https://docs.useautumn.com/quickstart/shadcn) to customize.

### [Integrate Pricing Logic](#integrate-pricing-logic)

Integrate your client and server pricing tiers logic with the following functions:

*   `check` to see if the customer is `allowed` to send a message.
*   `track` a usage event in Autumn (typically done server-side)
*   `customer` to display any relevant billing data in your UI (subscriptions, feature balances)

Server-side, you can access Autumn's functions through the `auth` object.

### [Additional Functions](#additional-functions)

#### [openBillingPortal()](#openbillingportal)

Opens a billing portal where the customer can update their payment method or cancel their plan.

    import { useCustomer } from "autumn-js/react";
    
    export default function BillingSettings() {
      const { openBillingPortal } = useCustomer();
    
      return (
        <button
          onClick={async () => {
            await openBillingPortal({
              returnUrl: "/settings/billing",
            });
          }}
        >
          Manage Billing
        </button>
      );
    }

#### [cancel()](#cancel)

Cancel a product or subscription.

    import { useCustomer } from "autumn-js/react";
    
    export default function CancelSubscription() {
      const { cancel } = useCustomer();
    
      return (
        <button
          onClick={async () => {
            await cancel({ productId: "pro" });
          }}
        >
          Cancel Subscription
        </button>
      );
    }

#### [Get invoice history](#get-invoice-history)

Pass in an `expand` param into `useCustomer` to get additional information. You can expand `invoices`, `trials_used`, `payment_method`, or `rewards`.

    import { useCustomer } from "autumn-js/react";
    
    export default function CustomerProfile() {
      const { customer } = useCustomer({ expand: ["invoices"] });
    
      return (
        <div>
          <h2>Customer Profile</h2>
          <p>Name: {customer?.name}</p>
          <p>Email: {customer?.email}</p>
          <p>Balance: {customer?.features.chat_messages?.balance}</p>
        </div>
      );
    }

[Edit on GitHub](https://github.com/better-auth/better-auth/blob/canary/docs/content/docs/plugins/autumn.mdx)</content>
</page>

<page>
  <title>Polar | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/polar</url>
  <content>[Polar](https://polar.sh/) is a developer first payment infrastructure. Out of the box it provides a lot of developer first integrations for payments, checkouts and more. This plugin helps you integrate Polar with Better Auth to make your auth + payments flow seamless.

This plugin is maintained by Polar team. For bugs, issues or feature requests, please visit the [Polar GitHub repo](https://github.com/polarsource/polar-adapters).

*   Checkout Integration
*   Customer Portal
*   Automatic Customer creation on signup
*   Event Ingestion & Customer Meters for flexible Usage Based Billing
*   Handle Polar Webhooks securely with signature verification
*   Reference System to associate purchases with organizations

    pnpm add better-auth @polar-sh/better-auth @polar-sh/sdk

Go to your Polar Organization Settings, and create an Organization Access Token. Add it to your environment.

    # .env
    POLAR_ACCESS_TOKEN=...

### [Configuring BetterAuth Server](#configuring-betterauth-server)

The Polar plugin comes with a handful additional plugins which adds functionality to your stack.

*   Checkout - Enables a seamless checkout integration
*   Portal - Makes it possible for your customers to manage their orders, subscriptions & granted benefits
*   Usage - Simple extension for listing customer meters & ingesting events for Usage Based Billing
*   Webhooks - Listen for relevant Polar webhooks

    import { betterAuth } from "better-auth";
    import { polar, checkout, portal, usage, webhooks } from "@polar-sh/better-auth";
    import { Polar } from "@polar-sh/sdk";
    
    const polarClient = new Polar({
        accessToken: process.env.POLAR_ACCESS_TOKEN,
        // Use 'sandbox' if you're using the Polar Sandbox environment
        // Remember that access tokens, products, etc. are completely separated between environments.
        // Access tokens obtained in Production are for instance not usable in the Sandbox environment.
        server: 'sandbox'
    });
    
    const auth = betterAuth({
        // ... Better Auth config
        plugins: [
            polar({
                client: polarClient,
                createCustomerOnSignUp: true,
                use: [
                    checkout({
                        products: [
                            {
                                productId: "123-456-789", // ID of Product from Polar Dashboard
                                slug: "pro" // Custom slug for easy reference in Checkout URL, e.g. /checkout/pro
                            }
                        ],
                        successUrl: "/success?checkout_id={CHECKOUT_ID}",
                        authenticatedUsersOnly: true
                    }),
                    portal(),
                    usage(),
                    webhooks({
                        secret: process.env.POLAR_WEBHOOK_SECRET,
                        onCustomerStateChanged: (payload) => // Triggered when anything regarding a customer changes
                        onOrderPaid: (payload) => // Triggered when an order was paid (purchase, subscription renewal, etc.)
                        ...  // Over 25 granular webhook handlers
                        onPayload: (payload) => // Catch-all for all events
                    })
                ],
            })
        ]
    });

### [Configuring BetterAuth Client](#configuring-betterauth-client)

You will be using the BetterAuth Client to interact with the Polar functionalities.

    import { createAuthClient } from "better-auth/react";
    import { polarClient } from "@polar-sh/better-auth";
    
    // This is all that is needed
    // All Polar plugins, etc. should be attached to the server-side BetterAuth config
    export const authClient = createAuthClient({
      plugins: [polarClient()],
    });

    import { betterAuth } from "better-auth";
    import {
      polar,
      checkout,
      portal,
      usage,
      webhooks,
    } from "@polar-sh/better-auth";
    import { Polar } from "@polar-sh/sdk";
    
    const polarClient = new Polar({
      accessToken: process.env.POLAR_ACCESS_TOKEN,
      // Use 'sandbox' if you're using the Polar Sandbox environment
      // Remember that access tokens, products, etc. are completely separated between environments.
      // Access tokens obtained in Production are for instance not usable in the Sandbox environment.
      server: "sandbox",
    });
    
    const auth = betterAuth({
      // ... Better Auth config
      plugins: [
        polar({
          client: polarClient,
          createCustomerOnSignUp: true,
          getCustomerCreateParams: ({ user }, request) => ({
            metadata: {
              myCustomProperty: 123,
            },
          }),
          use: [
            // This is where you add Polar plugins
          ],
        }),
      ],
    });

### [Required Options](#required-options)

*   `client`: Polar SDK client instance

### [Optional Options](#optional-options)

*   `createCustomerOnSignUp`: Automatically create a Polar customer when a user signs up
*   `getCustomerCreateParams`: Custom function to provide additional customer creation metadata

### [Customers](#customers)

When `createCustomerOnSignUp` is enabled, a new Polar Customer is automatically created when a new User is added in the Better-Auth Database.

All new customers are created with an associated `externalId`, which is the ID of your User in the Database. This allows us to skip any Polar to User mapping in your Database.

To support checkouts in your app, simply pass the Checkout plugin to the use-property.

    import { polar, checkout } from "@polar-sh/better-auth";
    
    const auth = betterAuth({
        // ... Better Auth config
        plugins: [
            polar({
                ...
                use: [
                    checkout({
                        // Optional field - will make it possible to pass a slug to checkout instead of Product ID
                        products: [ { productId: "123-456-789", slug: "pro" } ],
                        // Relative URL to return to when checkout is successfully completed
                        successUrl: "/success?checkout_id={CHECKOUT_ID}",
                        // Whether you want to allow unauthenticated checkout sessions or not
                        authenticatedUsersOnly: true
                    })
                ],
            })
        ]
    });

When checkouts are enabled, you're able to initialize Checkout Sessions using the checkout-method on the BetterAuth Client. This will redirect the user to the Product Checkout.

    await authClient.checkout({
      // Any Polar Product ID can be passed here
      products: ["e651f46d-ac20-4f26-b769-ad088b123df2"],
      // Or, if you setup "products" in the Checkout Config, you can pass the slug
      slug: "pro",
    });

Checkouts will automatically carry the authenticated User as the customer to the checkout. Email-address will be "locked-in".

If `authenticatedUsersOnly` is `false` - then it will be possible to trigger checkout sessions without any associated customer.

### [Organization Support](#organization-support)

This plugin supports the Organization plugin. If you pass the organization ID to the Checkout referenceId, you will be able to keep track of purchases made from organization members.

    const organizationId = (await authClient.organization.list())?.data?.[0]?.id,
    
    await authClient.checkout({
        // Any Polar Product ID can be passed here
        products: ["e651f46d-ac20-4f26-b769-ad088b123df2"],
        // Or, if you setup "products" in the Checkout Config, you can pass the slug
        slug: 'pro',
        // Reference ID will be saved as `referenceId` in the metadata of the checkout, order & subscription object
        referenceId: organizationId
    });

A plugin which enables customer management of their purchases, orders and subscriptions.

    import { polar, checkout, portal } from "@polar-sh/better-auth";
    
    const auth = betterAuth({
        // ... Better Auth config
        plugins: [
            polar({
                ...
                use: [
                    checkout(...),
                    portal()
                ],
            })
        ]
    });

The portal-plugin gives the BetterAuth Client a set of customer management methods, scoped under `authClient.customer`.

### [Customer Portal Management](#customer-portal-management)

The following method will redirect the user to the Polar Customer Portal, where they can see orders, purchases, subscriptions, benefits, etc.

    await authClient.customer.portal();

### [Customer State](#customer-state)

The portal plugin also adds a convenient state-method for retrieving the general Customer State.

    const { data: customerState } = await authClient.customer.state();

The customer state object contains:

*   All the data about the customer.
*   The list of their active subscriptions
    *   Note: This does not include subscriptions done by a parent organization. See the subscription list-method below for more information.
*   The list of their granted benefits.
*   The list of their active meters, with their current balance.

Thus, with that single object, you have all the required information to check if you should provision access to your service or not.

[You can learn more about the Polar Customer State in the Polar Docs](https://docs.polar.sh/integrate/customer-state).

### [Benefits, Orders & Subscriptions](#benefits-orders--subscriptions)

The portal plugin adds 3 convenient methods for listing benefits, orders & subscriptions relevant to the authenticated user/customer.

[All of these methods use the Polar CustomerPortal APIs](https://docs.polar.sh/api-reference/customer-portal)

#### [Benefits](#benefits)

This method only lists granted benefits for the authenticated user/customer.

    const { data: benefits } = await authClient.customer.benefits.list({
      query: {
        page: 1,
        limit: 10,
      },
    });

#### [Orders](#orders)

This method lists orders like purchases and subscription renewals for the authenticated user/customer.

    const { data: orders } = await authClient.customer.orders.list({
      query: {
        page: 1,
        limit: 10,
        productBillingType: "one_time", // or 'recurring'
      },
    });

#### [Subscriptions](#subscriptions)

This method lists the subscriptions associated with authenticated user/customer.

    const { data: subscriptions } = await authClient.customer.subscriptions.list({
      query: {
        page: 1,
        limit: 10,
        active: true,
      },
    });

**Important** - Organization Support

This will **not** return subscriptions made by a parent organization to the authenticated user.

However, you can pass a `referenceId` to this method. This will return all subscriptions associated with that referenceId instead of subscriptions associated with the user.

So in order to figure out if a user should have access, pass the user's organization ID to see if there is an active subscription for that organization.

    const organizationId = (await authClient.organization.list())?.data?.[0]?.id,
    
    const { data: subscriptions } = await authClient.customer.orders.list({
        query: {
    	    page: 1,
    		limit: 10,
    		active: true,
            referenceId: organizationId
        },
    });
    
    const userShouldHaveAccess = subscriptions.some(
        sub => // Your logic to check subscription product or whatever.
    )

A simple plugin for Usage Based Billing.

    import { polar, checkout, portal, usage } from "@polar-sh/better-auth";
    
    const auth = betterAuth({
        // ... Better Auth config
        plugins: [
            polar({
                ...
                use: [
                    checkout(...),
                    portal(),
                    usage()
                ],
            })
        ]
    });

### [Event Ingestion](#event-ingestion)

Polar's Usage Based Billing builds entirely on event ingestion. Ingest events from your application, create Meters to represent that usage, and add metered prices to Products to charge for it.

[Learn more about Usage Based Billing in the Polar Docs.](https://docs.polar.sh/features/usage-based-billing/introduction)

    const { data: ingested } = await authClient.usage.ingest({
      event: "file-uploads",
      metadata: {
        uploadedFiles: 12,
      },
    });

The authenticated user is automatically associated with the ingested event.

### [Customer Meters](#customer-meters)

A simple method for listing the authenticated user's Usage Meters, or as we call them, Customer Meters.

Customer Meter's contains all information about their consumption on your defined meters.

*   Customer Information
*   Meter Information
*   Customer Meter Information
    *   Consumed Units
    *   Credited Units
    *   Balance

    const { data: customerMeters } = await authClient.usage.meters.list({
      query: {
        page: 1,
        limit: 10,
      },
    });

The Webhooks plugin can be used to capture incoming events from your Polar organization.

    import { polar, webhooks } from "@polar-sh/better-auth";
    
    const auth = betterAuth({
        // ... Better Auth config
        plugins: [
            polar({
                ...
                use: [
                    webhooks({
                        secret: process.env.POLAR_WEBHOOK_SECRET,
                        onCustomerStateChanged: (payload) => // Triggered when anything regarding a customer changes
                        onOrderPaid: (payload) => // Triggered when an order was paid (purchase, subscription renewal, etc.)
                        ...  // Over 25 granular webhook handlers
                        onPayload: (payload) => // Catch-all for all events
                    })
                ],
            })
        ]
    });

Configure a Webhook endpoint in your Polar Organization Settings page. Webhook endpoint is configured at /polar/webhooks.

Add the secret to your environment.

    # .env
    POLAR_WEBHOOK_SECRET=...

The plugin supports handlers for all Polar webhook events:

*   `onPayload` - Catch-all handler for any incoming Webhook event
*   `onCheckoutCreated` - Triggered when a checkout is created
*   `onCheckoutUpdated` - Triggered when a checkout is updated
*   `onOrderCreated` - Triggered when an order is created
*   `onOrderPaid` - Triggered when an order is paid
*   `onOrderRefunded` - Triggered when an order is refunded
*   `onRefundCreated` - Triggered when a refund is created
*   `onRefundUpdated` - Triggered when a refund is updated
*   `onSubscriptionCreated` - Triggered when a subscription is created
*   `onSubscriptionUpdated` - Triggered when a subscription is updated
*   `onSubscriptionActive` - Triggered when a subscription becomes active
*   `onSubscriptionCanceled` - Triggered when a subscription is canceled
*   `onSubscriptionRevoked` - Triggered when a subscription is revoked
*   `onSubscriptionUncanceled` - Triggered when a subscription cancellation is reversed
*   `onProductCreated` - Triggered when a product is created
*   `onProductUpdated` - Triggered when a product is updated
*   `onOrganizationUpdated` - Triggered when an organization is updated
*   `onBenefitCreated` - Triggered when a benefit is created
*   `onBenefitUpdated` - Triggered when a benefit is updated
*   `onBenefitGrantCreated` - Triggered when a benefit grant is created
*   `onBenefitGrantUpdated` - Triggered when a benefit grant is updated
*   `onBenefitGrantRevoked` - Triggered when a benefit grant is revoked
*   `onCustomerCreated` - Triggered when a customer is created
*   `onCustomerUpdated` - Triggered when a customer is updated
*   `onCustomerDeleted` - Triggered when a customer is deleted
*   `onCustomerStateChanged` - Triggered when a customer is created</content>
</page>

<page>
  <title>Stripe | Better Auth</title>
  <url>https://www.better-auth.com/docs/plugins/stripe</url>
  <content>The Stripe plugin integrates Stripe's payment and subscription functionality with Better Auth. Since payment and authentication are often tightly coupled, this plugin simplifies the integration of Stripe into your application, handling customer creation, subscription management, and webhook processing.

*   Create Stripe Customers automatically when users sign up
*   Manage subscription plans and pricing
*   Process subscription lifecycle events (creation, updates, cancellations)
*   Handle Stripe webhooks securely with signature verification
*   Expose subscription data to your application
*   Support for trial periods and subscription upgrades
*   **Automatic trial abuse prevention** - Users can only get one trial per account across all plans
*   Flexible reference system to associate subscriptions with users or organizations
*   Team subscription support with seats management

### [Install the plugin](#install-the-plugin)

First, install the plugin:

If you're using a separate client and server setup, make sure to install the plugin in both parts of your project.

### [Install the Stripe SDK](#install-the-stripe-sdk)

Next, install the Stripe SDK on your server:

### [Add the plugin to your auth config](#add-the-plugin-to-your-auth-config)

auth.ts

    import { betterAuth } from "better-auth"
    import { stripe } from "@better-auth/stripe"
    import Stripe from "stripe"
    
    const stripeClient = new Stripe(process.env.STRIPE_SECRET_KEY!, {
        apiVersion: "2025-11-17.clover", // Latest API version as of Stripe SDK v20.0.0
    })
    
    export const auth = betterAuth({
        // ... your existing config
        plugins: [
            stripe({
                stripeClient,
                stripeWebhookSecret: process.env.STRIPE_WEBHOOK_SECRET!,
                createCustomerOnSignUp: true,
            })
        ]
    })

**Upgrading from Stripe v18?** Version 19 uses async webhook signature verification (`constructEventAsync`) which is handled internally by the plugin. No code changes required on your end!

### [Add the client plugin](#add-the-client-plugin)

auth-client.ts

    import { createAuthClient } from "better-auth/client"
    import { stripeClient } from "@better-auth/stripe/client"
    
    export const client = createAuthClient({
        // ... your existing config
        plugins: [
            stripeClient({
                subscription: true //if you want to enable subscription management
            })
        ]
    })

### [Migrate the database](#migrate-the-database)

Run the migration or generate the schema to add the necessary tables to the database.

See the [Schema](#schema) section to add the tables manually.

### [Set up Stripe webhooks](#set-up-stripe-webhooks)

Create a webhook endpoint in your Stripe dashboard pointing to:

    https://your-domain.com/api/auth/stripe/webhook

`/api/auth` is the default path for the auth server.

Make sure to select at least these events:

*   `checkout.session.completed`
*   `customer.subscription.updated`
*   `customer.subscription.deleted`

Save the webhook signing secret provided by Stripe and add it to your environment variables as `STRIPE_WEBHOOK_SECRET`.

### [Customer Management](#customer-management)

You can use this plugin solely for customer management without enabling subscriptions. This is useful if you just want to link Stripe customers to your users.

By default, when a user signs up, a Stripe customer is automatically created if you set `createCustomerOnSignUp: true`. This customer is linked to the user in your database. You can customize the customer creation process:

auth.ts

    stripe({
        // ... other options
        createCustomerOnSignUp: true,
        onCustomerCreate: async ({ stripeCustomer, user }, ctx) => {
            // Do something with the newly created customer
            console.log(`Customer ${stripeCustomer.id} created for user ${user.id}`);
        },
        getCustomerCreateParams: async (user, ctx) => {
            // Customize the Stripe customer creation parameters
            return {
                metadata: {
                    referralSource: user.metadata?.referralSource
                }
            };
        }
    })

### [Subscription Management](#subscription-management)

#### [Defining Plans](#defining-plans)

You can define your subscription plans either statically or dynamically:

auth.ts

    // Static plans
    subscription: {
        enabled: true,
        plans: [
            {
                name: "basic", // the name of the plan, it'll be automatically lower cased when stored in the database
                priceId: "price_1234567890", // the price ID from stripe
                annualDiscountPriceId: "price_1234567890", // (optional) the price ID for annual billing with a discount
                limits: {
                    projects: 5,
                    storage: 10
                }
            },
            {
                name: "pro",
                priceId: "price_0987654321",
                limits: {
                    projects: 20,
                    storage: 50
                },
                freeTrial: {
                    days: 14,
                }
            }
        ]
    }
    
    // Dynamic plans (fetched from database or API)
    subscription: {
        enabled: true,
        plans: async () => {
            const plans = await db.query("SELECT * FROM plans");
            return plans.map(plan => ({
                name: plan.name,
                priceId: plan.stripe_price_id,
                limits: JSON.parse(plan.limits)
            }));
        }
    }

see [plan configuration](#plan-configuration) for more.

#### [Creating a Subscription](#creating-a-subscription)

To create a subscription, use the `subscription.upgrade` method:

POST

/subscription/upgrade

    const { data, error } = await authClient.subscription.upgrade({    plan: "pro", // required    annual: true,    referenceId: "123",    subscriptionId: "sub_123",    metadata,    seats: 1,    successUrl, // required    cancelUrl, // required    returnUrl,    disableRedirect: true, // required});

| Prop | Description | Type |
| --- | --- | --- |
| `plan` | 
The name of the plan to upgrade to.

 | `string` |
| `annual?` | 

Whether to upgrade to an annual plan.

 | `boolean` |
| `referenceId?` | 

Reference id of the subscription to upgrade.

 | `string` |
| `subscriptionId?` | 

The id of the subscription to upgrade.

 | `string` |
| `metadata?` |  | `Record<string, any>` |
| `seats?` | 

Number of seats to upgrade to (if applicable).

 | `number` |
| `successUrl` | 

Callback URL to redirect back after successful subscription.

 | `string` |
| `cancelUrl` | 

If set, checkout shows a back button and customers will be directed here if they cancel payment.

 | `string` |
| `returnUrl?` | 

URL to take customers to when they click on the billing portal’s link to return to your website.

 | `string` |
| `disableRedirect` | 

Disable redirect after successful subscription.

 | `boolean` |

**Simple Example:**

client.ts

    await client.subscription.upgrade({
        plan: "pro",
        successUrl: "/dashboard",
        cancelUrl: "/pricing",
        annual: true, // Optional: upgrade to an annual plan
        referenceId: "org_123", // Optional: defaults to the current logged in user ID
        seats: 5 // Optional: for team plans
    });

This will create a Checkout Session and redirect the user to the Stripe Checkout page.

If the user already has an active subscription, you _must_ provide the `subscriptionId` parameter. Otherwise, the user will be subscribed to (and pay for) both plans.

> **Important:** The `successUrl` parameter will be internally modified to handle race conditions between checkout completion and webhook processing. The plugin creates an intermediate redirect that ensures subscription status is properly updated before redirecting to your success page.

    const { error } = await client.subscription.upgrade({
        plan: "pro",
        successUrl: "/dashboard",
        cancelUrl: "/pricing",
    });
    if(error) {
        alert(error.message);
    }

For each reference ID (user or organization), only one active or trialing subscription is supported at a time. The plugin doesn't currently support multiple concurrent active subscriptions for the same reference ID.

#### [Switching Plans](#switching-plans)

To switch a subscription to a different plan, use the `subscription.upgrade` method:

client.ts

    await client.subscription.upgrade({
        plan: "pro",
        successUrl: "/dashboard",
        cancelUrl: "/pricing",
        subscriptionId: "sub_123", // the Stripe subscription ID of the user's current plan
    });

This ensures that the user only pays for the new plan, and not both.

#### [Listing Active Subscriptions](#listing-active-subscriptions)

To get the user's active subscriptions:

    const { data: subscriptions, error } = await authClient.subscription.list({    query: {        referenceId: '123',    },});// get the active subscriptionconst activeSubscription = subscriptions.find(    sub => sub.status === "active" || sub.status === "trialing");// Check subscription limitsconst projectLimit = subscriptions?.limits?.projects || 0;

| Prop | Description | Type |
| --- | --- | --- |
| `referenceId?` | 
Reference id of the subscription to list.

 | `string` |

Make sure to provide `authorizeReference` in your plugin config to authorize the reference ID

auth.ts

    stripe({
        // ... other options
        subscription: {
            // ... other subscription options
            authorizeReference: async ({ user, session, referenceId, action }) => {
                if(action === "list-subscription") {
                    const org = await db.member.findFirst({
                        where: {
                            organizationId: referenceId,
                            userId: user.id
                        }   
                    });
                    return org?.role === "owner"
                }
                // Check if the user has permission to list subscriptions for this reference
                return true;
            }
        }
    })

#### [Canceling a Subscription](#canceling-a-subscription)

To cancel a subscription:

    const { data, error } = await authClient.subscription.cancel({    referenceId: 'org_123',    subscriptionId: 'sub_123',    returnUrl: '/account', // required});

| Prop | Description | Type |
| --- | --- | --- |
| `referenceId?` | 
Reference id of the subscription to cancel. Defaults to the userId.

 | `string` |
| `subscriptionId?` | 

The id of the subscription to cancel.

 | `string` |
| `returnUrl` | 

URL to take customers to when they click on the billing portal’s link to return to your website.

 | `string` |

This will redirect the user to the Stripe Billing Portal where they can cancel their subscription.

#### [Restoring a Canceled Subscription](#restoring-a-canceled-subscription)

If a user changes their mind after canceling a subscription (but before the subscription period ends), you can restore the subscription:

POST

/subscription/restore

    const { data, error } = await authClient.subscription.restore({    referenceId: '123',    subscriptionId: 'sub_123',});

| Prop | Description | Type |
| --- | --- | --- |
| `referenceId?` | 
Reference id of the subscription to restore. Defaults to the userId.

 | `string` |
| `subscriptionId?` | 

The id of the subscription to restore.

 | `string` |

This will reactivate a subscription that was previously set to cancel at the end of the billing period (`cancelAtPeriodEnd: true`). The subscription will continue to renew automatically.

> **Note:** This only works for subscriptions that are still active but marked to cancel at the end of the period. It cannot restore subscriptions that have already ended.

#### [Creating Billing Portal Sessions](#creating-billing-portal-sessions)

To create a [Stripe billing portal session](https://docs.stripe.com/api/customer_portal/sessions/create) where customers can manage their subscriptions, update payment methods, and view billing history:

POST

/subscription/billing-portal

    const { data, error } = await authClient.subscription.billingPortal({    locale,    referenceId: "123",    returnUrl,});

| Prop | Description | Type |
| --- | --- | --- |
| `locale?` | 
The IETF language tag of the locale customer portal is displayed in. If blank or auto, browser's locale is used.

 | `string` |
| `referenceId?` | 

Reference id of the subscription to upgrade.

 | `string` |
| `returnUrl?` | 

Return URL to redirect back after successful subscription.

 | `string` |

This endpoint creates a Stripe billing portal session and returns a URL in the response as `data.url`. You can redirect users to this URL to allow them to manage their subscription, payment methods, and billing history.

### [Reference System](#reference-system)

By default, subscriptions are associated with the user ID. However, you can use a custom reference ID to associate subscriptions with other entities, such as organizations:

client.ts

    // Create a subscription for an organization
    await client.subscription.upgrade({
        plan: "pro",
        referenceId: "org_123456",
        successUrl: "/dashboard",
        cancelUrl: "/pricing",
        seats: 5 // Number of seats for team plans
    });
    
    // List subscriptions for an organization
    const { data: subscriptions } = await client.subscription.list({
        query: {
            referenceId: "org_123456"
        }
    });

#### [Team Subscriptions with Seats](#team-subscriptions-with-seats)

For team or organization plans, you can specify the number of seats:

    await client.subscription.upgrade({
        plan: "team",
        referenceId: "org_123456",
        seats: 10, // 10 team members
        successUrl: "/org/billing/success",
        cancelUrl: "/org/billing"
    });

The `seats` parameter is passed to Stripe as the quantity for the subscription item. You can use this value in your application logic to limit the number of members in a team or organization.

To authorize reference IDs, implement the `authorizeReference` function:

auth.ts

    subscription: {
        // ... other options
        authorizeReference: async ({ user, session, referenceId, action }) => {
            // Check if the user has permission to manage subscriptions for this reference
            if (action === "upgrade-subscription" || action === "cancel-subscription" || action === "restore-subscription") {
                const org = await db.member.findFirst({
                    where: {
                        organizationId: referenceId,
                        userId: user.id
                    }   
                });
                return org?.role === "owner"
            }
            return true;
        }
    }

### [Webhook Handling](#webhook-handling)

The plugin automatically handles common webhook events:

*   `checkout.session.completed`: Updates subscription status after checkout
*   `customer.subscription.updated`: Updates subscription details when changed
*   `customer.subscription.deleted`: Marks subscription as canceled

You can also handle custom events:

auth.ts

    stripe({
        // ... other options
        onEvent: async (event) => {
            // Handle any Stripe event
            switch (event.type) {
                case "invoice.paid":
                    // Handle paid invoice
                    break;
                case "payment_intent.succeeded":
                    // Handle successful payment
                    break;
            }
        }
    })

### [Subscription Lifecycle Hooks](#subscription-lifecycle-hooks)

You can hook into various subscription lifecycle events:

auth.ts

    subscription: {
        // ... other options
        onSubscriptionComplete: async ({ event, subscription, stripeSubscription, plan }) => {
            // Called when a subscription is successfully created
            await sendWelcomeEmail(subscription.referenceId, plan.name);
        },
        onSubscriptionUpdate: async ({ event, subscription }) => {
            // Called when a subscription is updated
            console.log(`Subscription ${subscription.id} updated`);
        },
        onSubscriptionCancel: async ({ event, subscription, stripeSubscription, cancellationDetails }) => {
            // Called when a subscription is canceled
            await sendCancellationEmail(subscription.referenceId);
        },
        onSubscriptionDeleted: async ({ event, subscription, stripeSubscription }) => {
            // Called when a subscription is deleted
            console.log(`Subscription ${subscription.id} deleted`);
        }
    }

### [Trial Periods](#trial-periods)

You can configure trial periods for your plans:

auth.ts

    {
        name: "pro",
        priceId: "price_0987654321",
        freeTrial: {
            days: 14,
            onTrialStart: async (subscription) => {
                // Called when a trial starts
                await sendTrialStartEmail(subscription.referenceId);
            },
            onTrialEnd: async ({ subscription }, ctx) => {
                // Called when a trial ends
                await sendTrialEndEmail(subscription.referenceId);
            },
            onTrialExpired: async (subscription, ctx) => {
                // Called when a trial expires without conversion
                await sendTrialExpiredEmail(subscription.referenceId);
            }
        }
    }

The Stripe plugin adds the following tables to your database:

### [User](#user)

Table Name: `user`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| stripeCustomerId | string |  | The Stripe customer ID |

### [Subscription](#subscription)

Table Name: `subscription`

| Field Name | Type | Key | Description |
| --- | --- | --- | --- |
| id | string |  | Unique identifier for each subscription |
| plan | string | \- | The name of the subscription plan |
| referenceId | string | \- | The ID this subscription is associated with (user ID by default). This should NOT be a unique field in your database, as it must allow users to resubscribe after a cancellation. |
| stripeCustomerId | string |  | The Stripe customer ID |
| stripeSubscriptionId | string |  | The Stripe subscription ID |
| status | string | \- | The status of the subscription (active, canceled, etc.) |
| periodStart | Date |  | Start date of the current billing period |
| periodEnd | Date |  | End date of the current billing period |
| cancelAtPeriodEnd | boolean |  | Whether the subscription will be canceled at the end of the period |
| seats | number |  | Number of seats for team plans |
| trialStart | Date |  | Start date of the trial period |
| trialEnd | Date |  | End date of the trial period |

### [Customizing the Schema](#customizing-the-schema)

To change the schema table names or fields, you can pass a `schema` option to the Stripe plugin:

auth.ts

    stripe({
        // ... other options
        schema: {
            subscription: {
                modelName: "stripeSubscriptions", // map the subscription table to stripeSubscriptions
                fields: {
                    plan: "planName" // map the plan field to planName
                }
            }
        }
    })

### [Main Options](#main-options)

**stripeClient**: `Stripe` - The Stripe client instance. Required.

**stripeWebhookSecret**: `string` - The webhook signing secret from Stripe. Required.

**createCustomerOnSignUp**: `boolean` - Whether to automatically create a Stripe customer when a user signs up. Default: `false`.

**onCustomerCreate**: `(data: { stripeCustomer: Stripe.Customer, user: User }, ctx: GenericEndpointContext) => Promise<void>` - A function called after a customer is created.

**getCustomerCreateParams**: `(user: User, ctx: GenericEndpointContext) => Promise<{}>` - A function to customize the Stripe customer creation parameters.

**onEvent**: `(event: Stripe.Event) => Promise<void>` - A function called for any Stripe webhook event.

### [Subscription Options](#subscription-options)

**enabled**: `boolean` - Whether to enable subscription functionality. Required.

**plans**: `Plan[] | (() => Promise<Plan[]>)` - An array of subscription plans or a function that returns plans. Required if subscriptions are enabled.

**requireEmailVerification**: `boolean` - Whether to require email verification before allowing subscription upgrades. Default: `false`.

**authorizeReference**: `(data: { user: User, session: Session, referenceId: string, action: "upgrade-subscription" | "list-subscription" | "cancel-subscription" | "restore-subscription"}, ctx: GenericEndpointContext) => Promise<boolean>` - A function to authorize reference IDs.

### [Plan Configuration](#plan-configuration)

Each plan can have the following properties:

**name**: `string` - The name of the plan. Required.

**priceId**: `string` - The Stripe price ID. Required unless using `lookupKey`.

**lookupKey**: `string` - The Stripe price lookup key. Alternative to `priceId`.

**annualDiscountPriceId**: `string` - A price ID for annual billing.

**annualDiscountLookupKey**: `string` - The Stripe price lookup key for annual billing. Alternative to `annualDiscountPriceId`.

**limits**: `Record<string, unknown>` - Limits associated with the plan (e.g., `{ projects: 10, storage: 5 }`). Useful when you want to define plan-specific metadata.

**group**: `string` - A group name for the plan, useful for categorizing plans.

**freeTrial**: Object containing trial configuration:

*   **days**: `number` - Number of trial days.
*   **onTrialStart**: `(subscription: Subscription) => Promise<void>` - Called when a trial starts.
*   **onTrialEnd**: `(data: { subscription: Subscription }, ctx: GenericEndpointContext) => Promise<void>` - Called when a trial ends.
*   **onTrialExpired**: `(subscription: Subscription, ctx: GenericEndpointContext) => Promise<void>` - Called when a trial expires without conversion.

### [Using with Organizations](#using-with-organizations)

The Stripe plugin works well with the organization plugin. You can associate subscriptions with organizations instead of individual users:

client.ts

    // Get the active organization
    const { data: activeOrg } = client.useActiveOrganization();
    
    // Create a subscription for the organization
    await client.subscription.upgrade({
        plan: "team",
        referenceId: activeOrg.id,
        seats: 10,
        annual: true, // upgrade to an annual plan (optional)
        successUrl: "/org/billing/success",
        cancelUrl: "/org/billing"
    });

Make sure to implement the `authorizeReference` function to verify that the user has permission to manage subscriptions for the organization:

auth.ts

    subscription: {
        // ... other subscription options
        authorizeReference: async ({ user, referenceId, action }) => {
            const member = await db.members.findFirst({
                where: {
                    userId: user.id,
                    organizationId: referenceId
                }
            });
            
            return member?.role === "owner" || member?.role === "admin";
        }
    }

### [Custom Checkout Session Parameters](#custom-checkout-session-parameters)

You can customize the Stripe Checkout session with additional parameters:

auth.ts

    getCheckoutSessionParams: async ({ user, session, plan, subscription }, ctx) => {
        return {
            params: {
                allow_promotion_codes: true,
                tax_id_collection: {
                    enabled: true
                },
                billing_address_collection: "required",
                custom_text: {
                    submit: {
                        message: "We'll start your subscription right away"
                    }
                },
                metadata: {
                    planType: "business",
                    referralCode: user.metadata?.referralCode
                }
            },
            options: {
                idempotencyKey: `sub_${user.id}_${plan.name}_${Date.now()}`
            }
        };
    }

### [Tax Collection](#tax-collection)

To collect tax IDs from the customer, set `tax_id_collection` to true:

auth.ts

    subscription: {
        // ... other options
        getCheckoutSessionParams: async ({ user, session, plan, subscription }, ctx) => {
            return {
                params: {
                    tax_id_collection: {
                        enabled: true
                    }
                }
            };
        }
    }

### [Automatic Tax Calculation](#automatic-tax-calculation)

To enable automatic tax calculation using the customer's location, set `automatic_tax` to true. Enabling this parameter causes Checkout to collect any billing address information necessary for tax calculation. You need to have tax registration setup and configured in the Stripe dashboard first for this to work.

auth.ts

    subscription: {
        // ... other options
        getCheckoutSessionParams: async ({ user, session, plan, subscription }, ctx) => {
            return {
                params: {
                    automatic_tax: {
                        enabled: true
                    }
                }
            };
        }
    }

### [Trial Period Management](#trial-period-management)

The Stripe plugin automatically prevents users from getting multiple free trials. Once a user has used a trial period (regardless of which plan), they will not be eligible for additional trials on any plan.

**How it works:**

*   The system tracks trial usage across all plans for each user
*   When a user subscribes to a plan with a trial, the system checks their subscription history
*   If the user has ever had a trial (indicated by `trialStart`/`trialEnd` fields or `trialing` status), no new trial will be offered
*   This prevents abuse where users cancel subscriptions and resubscribe to get multiple free trials

**Example scenario:**

1.  User subscribes to "Starter" plan with 7-day trial
2.  User cancels the subscription after the trial
3.  User tries to subscribe to "Premium" plan - no trial will be offered
4.  User will be charged immediately for the Premium plan

This behavior is automatic and requires no additional configuration. The trial eligibility is determined at the time of subscription creation and cannot be overridden through configuration.

### [Webhook Issues](#webhook-issues)

If webhooks aren't being processed correctly:

1.  Check that your webhook URL is correctly configured in the Stripe dashboard
2.  Verify that the webhook signing secret is correct
3.  Ensure you've selected all the necessary events in the Stripe dashboard
4.  Check your server logs for any errors during webhook processing

### [Subscription Status Issues](#subscription-status-issues)

If subscription statuses aren't updating correctly:

1.  Make sure the webhook events are being received and processed
2.  Check that the `stripeCustomerId` and `stripeSubscriptionId` fields are correctly populated
3.  Verify that the reference IDs match between your application and Stripe

### [Testing Webhooks Locally](#testing-webhooks-locally)

For local development, you can use the Stripe CLI to forward webhooks to your local environment:

    stripe listen --forward-to localhost:3000/api/auth/stripe/webhook

This will provide you with a webhook signing secret that you can use in your local environment.</content>
</page>