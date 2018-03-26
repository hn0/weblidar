
// TODO: pass values from wrapper!
__kernel void euclid_dist(__global float* inx,
                          __global float* iny,
                          __global float* inz,
                          __global float* dist,
                          __global float* angle) {

    float3 origin, zaxis;

    const uint start = get_local_id(0);
    const uint sz    = get_global_size(0);
    const uint off   = get_local_size(0);

    origin = (float3)( 0.0 );
    zaxis  = (float3)( 1.0, 1.0, 1.0 );

    for( int i = start; i < sz; i += off ){
        dist[i] = fabs( distance( origin, (float3)( inx[i], iny[i], inz[i]) ) );
        angle[i] = dot( zaxis, (float3)(inx[i], iny[i], inz[i]));
    }
}