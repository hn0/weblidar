
__kernel void euclid_dist(__global float* inx,
                          __global float* iny,
                          __global float* inz,
                          __global float* dist) {
    
    const uint start = get_local_id(0);
    const uint sz    = get_global_size(0);
    const uint off   = get_local_size(0);


    for( int i = start; i < sz; i += off ){
        // float cnt = 0;
        // for( int j = i-1; j > 0; j -= off ){
        //     cnt++;
        // }
        // out[i] = in[i] * 100 + cnt;

        // distance from origin
        dist[i] = sqrt( inx[i] + iny[i] + inz[i] );
    }
}