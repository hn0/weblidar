

function LidarCanvas(app, stream, container, info_data){

    this.canvas = document.createElement( 'canvas' );
    container.appendChild( this.canvas );

    if( !this.init_canvas() ){
        console.warn( 'Something about missing gl canvas should be written!' );
        return;
    }

    stream.on( 'pointbatch', function(cnt){
        console.log( 'got the pts batch', cnt, this );
    }.bind( this ) );

    if( this.init_shaders() ){
        this.init_buffers();
        stream.start_streaming();
    }

};

LidarCanvas.prototype.init_canvas = function()
{
    var support = false;
    try {
        let gl = this.canvas.getContext( 'webgl' );
        gl.viewportWidth  = this.canvas.width;
        gl.viewportHeight = this.canvas.height;
        support = true;

        this.shaderp = gl.createProgram();
    }
    catch (e) {}
    return support;
};

LidarCanvas.prototype.init_buffers = function()
{
    let gl = this.canvas.getContext( 'webgl' );
    this.shaderp.positionLocation = gl.getAttribLocation( this.shaderp, 'Pos' );
    gl.enableVertexAttribArray( this.shaderp.positionLocation );

    this.shaderp.persp = gl.getUniformLocation( this.shaderp, 'u_persp' );
    this.shaderp.view  = gl.getUniformLocation( this.shaderp, 'u_modelview' );

    let sample = [
         0.0,  1.0, 0.0,
        -1.0, -1.0, 0.0,
         1.0, -1.0, 0.0
    ];

    let vecPos = gl.createBuffer();
    vecPos.itemSize = 3;
    vecPos.numItems = 3;
    gl.bindBuffer( gl.ARRAY_BUFFER, vecPos );
    gl.bufferData( gl.ARRAY_BUFFER, new Float32Array( sample ), gl.STATIC_DRAW );
    gl.vertexAttribPointer( this.shaderp.positionLocation, vecPos.itemSize, gl.FLOAT, false, 0, 0 );


    // for now draw sample image!
    gl.viewport( 0, 0, gl.viewportWidth, gl.viewportHeight );
    gl.clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT );

    // issue with mat4 is the library!
    let mvmat = mat4.create();
    let pmat  = mat4.create();
    mat4.identity( mvmat ); // wth, what this means?!
    mat4.translate( mvmat, mvmat, [0.0, 0.0, -3.5] );
    mat4.perspective( pmat, 45, gl.viewportWidth / gl.viewportHeight, .1, 100.0 );

    gl.uniformMatrix4fv( this.shaderp.persp, false, pmat );
    gl.uniformMatrix4fv( this.shaderp.view, false, mvmat );

    gl.drawArrays( gl.POINTS, 0, vecPos.numItems );
};

LidarCanvas.prototype.init_shaders = function()
{

    let gl = this.canvas.getContext( 'webgl' );
    let vertexs   = this.get_shader( 'shader-vertex', gl );
    let fragments = this.get_shader( 'shader-fragment', gl );

    if( !(vertexs || fragments) ){
        console.error( 'Could not load the shader programs!' );
        return false;
    }

    gl.attachShader( this.shaderp, vertexs );
    gl.attachShader( this.shaderp, fragments );
    gl.linkProgram( this.shaderp );

    if( !gl.getProgramParameter( this.shaderp, gl.LINK_STATUS ) ){
        console.error( 'Could not init the program!' );
        return false;
    }

    gl.useProgram( this.shaderp );
    return true;

};

LidarCanvas.prototype.get_shader = function(name, gl)
{
    let shader = null;
    let script = document.getElementById( name );

    if( script ){
        switch( script.type ){
            case 'x-shader/x-fragment':
                shader = gl.createShader( gl.FRAGMENT_SHADER );
                break;
            case 'x-shader/x-vertex':
                shader = gl.createShader( gl.VERTEX_SHADER );
                break;
        }

        if( shader ){
            gl.shaderSource( shader, script.text );
            gl.compileShader( shader );

            if( !gl.getShaderParameter( shader, gl.COMPILE_STATUS ) ){
                console.error( 'Shader compile error!' );
                console.error( gl.getShaderInfoLog( shader ) );
                shader = null;
            }
        }
    }

    return shader;
};