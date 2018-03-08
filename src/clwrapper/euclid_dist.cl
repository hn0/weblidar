
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
        // distance from origin
        // dist[i] = sqrt( pown(inx[i], 2) + pown(iny[i], 2) + pown(inz[i], 2) );

        // As of start, a crude catergorisation of the points is needed!
        //  having only magnitude and angle in respect to z axis should be enough to perform the tests
        dist[i] = distance( origin, (float3)( inx[i], iny[i], inz[i]) );
        // angle[i] = dot( zaxis, (float3)( inx[i], iny[i], inz[i]) );
        angle[i] = 10;
    }
}