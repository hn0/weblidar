
__kernel void square_test(__global float* in,
                          __global float* out) {   

   int i = get_global_id(0);
   out[i] = in[i] * in[i];
}