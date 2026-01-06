#ifdef __cplusplus
extern "C" {
#endif

__global__ void hello(char *message) {
    printf("%s\n", message);
}

#ifdef __cplusplus
}
#endif