/*
 * h2core.h - C API for h2.core HTTPS VPN
 *
 * This header defines the C-compatible API for integrating h2.core
 * into native applications via FFI (Foreign Function Interface).
 *
 * Usage:
 *   H2Instance inst = h2_create(config_json);
 *   h2_start(inst);
 *   // ... use VPN ...
 *   h2_stop(inst);
 *   h2_destroy(inst);
 */

#ifndef H2CORE_H
#define H2CORE_H

#ifdef __cplusplus
extern "C" {
#endif

/* Opaque handle to h2.core instance */
typedef void* H2Instance;

/* Error codes */
#define H2_OK              0
#define H2_ERR_NULL_PTR   -1
#define H2_ERR_INVALID    -2
#define H2_ERR_INIT       -3
#define H2_ERR_START      -4
#define H2_ERR_STOPPED    -5
#define H2_ERR_CONFIG     -6
#define H2_ERR_NETWORK    -7

/*
 * Server Mode API
 * ---------------
 * For running h2.core as a VPN server.
 */

/*
 * h2_create - Create a new h2.core server instance
 * @config_json: JSON configuration string (xray-compatible format)
 * @return: Instance handle, or NULL on error
 */
H2Instance h2_create(const char* config_json);

/*
 * h2_start - Start the VPN server
 * @instance: Instance handle from h2_create
 * @return: H2_OK on success, error code on failure
 */
int h2_start(H2Instance instance);

/*
 * h2_stop - Stop the VPN server
 * @instance: Instance handle
 * @return: H2_OK on success, error code on failure
 */
int h2_stop(H2Instance instance);

/*
 * h2_destroy - Destroy instance and free resources
 * @instance: Instance handle
 */
void h2_destroy(H2Instance instance);

/*
 * Client Mode API
 * ---------------
 * For connecting to an h2.core VPN server.
 * Creates a local SOCKS5 proxy that tunnels traffic through the VPN.
 */

/*
 * h2_client_create - Create a VPN client instance
 * @server_addr: VPN server address (e.g., "vpn.example.com:443")
 * @crypto_provider: Crypto provider name (e.g., "us", "ua", "cn")
 * @return: Instance handle, or NULL on error
 */
H2Instance h2_client_create(const char* server_addr, const char* crypto_provider);

/*
 * h2_client_connect - Connect to VPN server and start local SOCKS5 proxy
 * @instance: Client instance handle
 * @return: H2_OK on success, error code on failure
 */
int h2_client_connect(H2Instance instance);

/*
 * h2_client_disconnect - Disconnect from VPN server
 * @instance: Client instance handle
 * @return: H2_OK on success, error code on failure
 */
int h2_client_disconnect(H2Instance instance);

/*
 * h2_client_get_socks_port - Get the local SOCKS5 proxy port
 * @instance: Client instance handle
 * @return: Port number (> 0), or error code (< 0)
 */
int h2_client_get_socks_port(H2Instance instance);

/*
 * Information API
 * ---------------
 */

/*
 * h2_version - Get h2.core version string
 * @return: Version string (caller must NOT free)
 */
const char* h2_version(void);

/*
 * h2_get_stats - Get runtime statistics as JSON
 * @instance: Instance handle
 * @return: JSON string (caller must free with h2_free_string)
 */
const char* h2_get_stats(H2Instance instance);

/*
 * h2_is_running - Check if instance is running
 * @instance: Instance handle
 * @return: 1 if running, 0 if not, negative on error
 */
int h2_is_running(H2Instance instance);

/*
 * h2_get_last_error - Get last error message
 * @return: Error message string (caller must NOT free)
 */
const char* h2_get_last_error(void);

/*
 * Memory Management
 * -----------------
 */

/*
 * h2_free_string - Free a string returned by h2_get_stats
 * @str: String to free
 */
void h2_free_string(const char* str);

#ifdef __cplusplus
}
#endif

#endif /* H2CORE_H */
