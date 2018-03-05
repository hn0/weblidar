
__kernel void square_test(__global float4* in,
                          __global float4* out) {   
   int i = get_global_id(0);
   out[i] = in[i] * in[i];
}