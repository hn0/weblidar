__kernel void euclid_dist(__global float* inx,
                          __global float* iny,
                          __global float* inz,
                          __global float* cat) {
    
    // we need global variable to write data in!


    float3 origin, zaxis;

    const uint start = get_local_id(0);
    const uint sz    = get_global_size(0);
    const uint off   = get_local_size(0);

    origin = (float3)( 0.0 );
    zaxis  = (float3)( 1.0, 1.0, 1.0 );

    for( int i = start; i < sz; i += off ){

        // As of start, a crude catergorisation of the points is needed!
        //  having only magnitude and angle in respect to z axis should be enough to perform the tests
        // cat[i] = distance( origin, (float3)( inx[i], iny[i], inz[i]) );
        // cat[i] = start;
        cat[0]++;
    }
}