
// https://medium.com/social-tables-tech/hello-world-webgl-79f430446b5c
// http://learningwebgl.com/blog/?p=28

function LidarCanvas(app, stream, container, info_data){

    this.canvas = document.createElement( 'canvas' );

    if( !this.init_canvas ){
        console.warn( 'Something about missing gl canvas should be written!' );
        return;
    }

    if( this.init_shaders() ){
        stream.start_streaming();
    }

    // app.stream.init(this);
	console.log( 'stream should be started here!' );

};

LidarCanvas.prototype.init_canvas = function()
{
    var support = false;
    try {
        let gl = this.canvas.getContext( 'webgl' );
        gl.viewportWidth  = this.canvas.width;
        gl.viewportHeight = this.canvas.height;
        support = true;
    }
    catch (e) {}
    return support;
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

    let shaderp = gl.createProgram();
    gl.attachShader( shaderp, vertexs );
    gl.attachShader( shaderp, fragments );
    gl.linkProgram( shaderp );

    if( !gl.getProgramParameter( shaderp, gl.LINK_STATUS ) ){
        console.error( 'Could not init the program!' );
        return false;
    }

    gl.useProgram( shaderp );
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