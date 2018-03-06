
__kernel void euclid_dist(__global float* in,
                          __global float* out) {
    const uint id = get_local_id(0);
    const uint sz = get_local_size(0);
    ; const uint off = get_global_offset(0);
    out[id] = (float)(sz);
}