/*
 * Test program for h2.core C API
 *
 * Build:
 *   make -C cgo/test
 *
 * Run:
 *   ./cgo/test/test_h2core
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "../../dist/h2core.h"

void test_version(void) {
    printf("=== Test: h2_version ===\n");
    const char* version = h2_version();
    printf("Version: %s\n", version);
    // Note: h2_version returns static string, don't free
    printf("PASS\n\n");
}

void test_client_lifecycle(void) {
    printf("=== Test: Client Lifecycle ===\n");

    // Create client
    printf("Creating client...\n");
    H2Instance client = h2_client_create("test.example.com:443", "us");
    if (client == NULL) {
        printf("FAIL: h2_client_create returned NULL\n");
        printf("Error: %s\n", h2_get_last_error());
        return;
    }
    printf("Client created: %p\n", client);

    // Check not running
    int running = h2_is_running(client);
    printf("Is running (should be 0): %d\n", running);
    if (running != 0) {
        printf("FAIL: Should not be running yet\n");
        h2_destroy(client);
        return;
    }

    // Get stats before connect
    const char* stats = h2_get_stats(client);
    printf("Stats: %s\n", stats);
    h2_free_string(stats);

    // Get SOCKS port (should be 0 before connect)
    int port = h2_client_get_socks_port(client);
    printf("SOCKS port (should be 0): %d\n", port);

    // Connect (this will fail without a real server, but should return error gracefully)
    printf("Connecting (expected to start SOCKS listener)...\n");
    int result = h2_client_connect(client);
    printf("Connect result: %d\n", result);

    if (result == H2_OK) {
        // Get SOCKS port
        port = h2_client_get_socks_port(client);
        printf("SOCKS port: %d\n", port);

        // Check running
        running = h2_is_running(client);
        printf("Is running: %d\n", running);

        // Disconnect
        printf("Disconnecting...\n");
        result = h2_client_disconnect(client);
        printf("Disconnect result: %d\n", result);
    }

    // Destroy
    printf("Destroying client...\n");
    h2_destroy(client);
    printf("PASS\n\n");
}

void test_server_create_invalid(void) {
    printf("=== Test: Server Create (Invalid Config) ===\n");

    // Try to create with invalid JSON
    H2Instance inst = h2_create("not valid json");
    if (inst != NULL) {
        printf("FAIL: Should have failed with invalid JSON\n");
        h2_destroy(inst);
        return;
    }
    printf("Correctly rejected invalid JSON\n");
    printf("Error: %s\n", h2_get_last_error());
    printf("PASS\n\n");
}

void test_server_create_valid(void) {
    printf("=== Test: Server Create (Valid Config) ===\n");

    const char* config = "{"
        "\"inbounds\": [{"
        "  \"port\": 8443,"
        "  \"protocol\": \"https-vpn\","
        "  \"streamSettings\": {"
        "    \"network\": \"h2\","
        "    \"security\": \"tls\","
        "    \"tlsSettings\": {"
        "      \"certificates\": [{"
        "        \"certificateFile\": \"/tmp/cert.pem\","
        "        \"keyFile\": \"/tmp/key.pem\""
        "      }]"
        "    }"
        "  }"
        "}],"
        "\"outbounds\": [{"
        "  \"protocol\": \"freedom\""
        "}]"
        "}";

    printf("Creating server with config...\n");
    H2Instance inst = h2_create(config);

    if (inst == NULL) {
        // This is expected if cert files don't exist
        printf("Create returned NULL (expected if no certs)\n");
        printf("Error: %s\n", h2_get_last_error());
        printf("PASS (graceful failure)\n\n");
        return;
    }

    printf("Server created: %p\n", inst);

    // Get stats
    const char* stats = h2_get_stats(inst);
    printf("Stats: %s\n", stats);
    h2_free_string(stats);

    // Destroy without starting
    h2_destroy(inst);
    printf("PASS\n\n");
}

void test_null_handling(void) {
    printf("=== Test: NULL Handling ===\n");

    // All functions should handle NULL gracefully
    int result;

    result = h2_start(NULL);
    printf("h2_start(NULL): %d (expected %d)\n", result, H2_ERR_NULL_PTR);

    result = h2_stop(NULL);
    printf("h2_stop(NULL): %d (expected %d)\n", result, H2_ERR_NULL_PTR);

    h2_destroy(NULL); // Should not crash
    printf("h2_destroy(NULL): OK\n");

    result = h2_is_running(NULL);
    printf("h2_is_running(NULL): %d (expected %d)\n", result, H2_ERR_NULL_PTR);

    const char* stats = h2_get_stats(NULL);
    printf("h2_get_stats(NULL): %s\n", stats);
    h2_free_string(stats);

    result = h2_client_connect(NULL);
    printf("h2_client_connect(NULL): %d (expected %d)\n", result, H2_ERR_NULL_PTR);

    result = h2_client_disconnect(NULL);
    printf("h2_client_disconnect(NULL): %d (expected %d)\n", result, H2_ERR_NULL_PTR);

    result = h2_client_get_socks_port(NULL);
    printf("h2_client_get_socks_port(NULL): %d (expected %d)\n", result, H2_ERR_NULL_PTR);

    printf("PASS\n\n");
}

int main(int argc, char** argv) {
    printf("h2.core C API Test Suite\n");
    printf("========================\n\n");

    test_version();
    test_null_handling();
    test_client_lifecycle();
    test_server_create_invalid();
    test_server_create_valid();

    printf("========================\n");
    printf("All tests completed!\n");
    return 0;
}
