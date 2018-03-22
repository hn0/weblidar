

function LidarCanvas(stream, info){

    this.info = info;
    var container = document.getElementById( 'cancontainer' );
    if( !container ){
        return;
    }
    this.viewport = [0, 0, -4];
    this.canvas = document.createElement( 'canvas' );
    container.appendChild( this.canvas );

    if( !this.init_canvas() ){
        console.warn( 'Something about missing gl canvas should be written!' );
        return;
    }

    stream.on( 'pointbatch', function(){
        // this.draw_buffers( stream.coords, stream.points );
    }.bind( this ) );

    stream.on( 'done', function() {
        this.init_buffers( stream.coords, stream.points );
    }.bind( this ));

    if( this.init_shaders() ){
        stream.start_streaming();
        info.start_progress();
    }

    let self   = this;
    let mcords = [0, 0];
    let delta  = [0, 0];
    this.canvas.addEventListener( 'mousemove', function(evt) {
        if( evt.buttons == 1 ){
            delta = [evt.clientX, evt.clientY].map( function(v, i) { return mcords[i] - v; });

            // TODO: movement scale is probably off
            setTimeout(function(){
                delta.forEach( function(_, i){
                    let dir = (i == 0) ? -1 : 1;
                    self.viewport[i] += delta[i] * .01 * dir;
                });
                self.set_viewport();
                delta = [0, 0];
            }, 100);
        }
        mcords = [evt.clientX, evt.clientY];
    });

    //TODO: non standard feature, fix this!
    this.canvas.addEventListener( 'mousewheel', function(evt){
        const dir = Math.sign( evt.deltaY );
        self.viewport[2] += .2 * dir;
        self.set_viewport();
    });

};

LidarCanvas.prototype.set_viewport = function()
{
    let gl = this.canvas.getContext( 'webgl' );
    let mvmat = mat4.create();
    mat4.translate( mvmat, mvmat, this.viewport );
    this.shaderp.view  = gl.getUniformLocation( this.shaderp, 'u_modelview' );
    gl.uniformMatrix4fv( this.shaderp.view, false, mvmat );
    this.draw();
    this.info.set_viewoport( this.viewport );
};

LidarCanvas.prototype.draw = function()
{
    let gl = this.canvas.getContext( 'webgl' );

    gl.viewport( 0, 0, gl.viewportWidth, gl.viewportHeight );
    gl.clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT );

    gl.drawArrays( gl.POINTS, 0, this.vecPos.numItems );
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

LidarCanvas.prototype.init_buffers = function(pts, len)
{
    // ok, need to pull buffer data out!
    let gl = this.canvas.getContext( 'webgl' );
    this.shaderp.positionLocation = gl.getAttribLocation( this.shaderp, 'Pos' );
    gl.enableVertexAttribArray( this.shaderp.positionLocation );

    this.shaderp.persp = gl.getUniformLocation( this.shaderp, 'u_persp' );
    let pmat  = mat4.create();
    mat4.perspective( pmat, 45, gl.viewportWidth / gl.viewportHeight, .1, 100.0 );
    gl.uniformMatrix4fv( this.shaderp.persp, false, pmat );

    // let sample = [
    //      0.0,  1.0, 0.0,
    //     -1.0, -1.0, 0.0,
    //      1.0, -1.0, 0.0
    // ];

    this.vecPos = gl.createBuffer();
    this.vecPos.itemSize = 3;
    this.vecPos.numItems = len;
    gl.bindBuffer( gl.ARRAY_BUFFER, this.vecPos );
    // gl.bufferData( gl.ARRAY_BUFFER, new Float32Array( sample ), gl.STATIC_DRAW );
    gl.bufferData( gl.ARRAY_BUFFER, pts, gl.STATIC_DRAW );
    gl.vertexAttribPointer( this.shaderp.positionLocation, this.vecPos.itemSize, gl.FLOAT, false, 0, 0 );
    
    this.set_viewport();
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