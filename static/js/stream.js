

function Stream (){
    this.request = 0;
    this.points  = [];

    this.streams = [0, 1, 2];
    this.streams.forEach( this.init_stream.bind( this ) );
}

Stream.prototype.init_stream = function(reqid)
{
    this.get( reqid )
        .then( function( data ){
            if( this.parse_response( data ) ){
                this.init_stream( reqid + this.streams.length );
            }
        }.bind( this ));
};

Stream.prototype.parse_response = function(data)
{
    if( data.byteLength > 4 ){
        // first byte always is length of response
        var len = new Int32Array( data.slice( 0, 4 ) );
        var pts = new Float32Array( data.slice( 4 ) );
        if( !len.length || !len[0] ){
            console.log( 'stream done' );
            return false;
        }
        
        console.log( 'Got n pts for processing:', len[0] );
        for( i=0; i < len[0]; i++ ){
            console.log( 'pt', i, 'data', pts[i] )
        }

        return true;
    }
    return true;
};

Stream.prototype.get = function( reqid ){

    // we will need a response size!?
    return new Promise( (success, error) => {
            
        var req = new XMLHttpRequest();
        req.responseType = 'arraybuffer';
        req.onload  = () => success( req.response );
        req.onerror = () => error( null );

        try{
            // maybe some kind of clientid will be needed!
            req.open( "GET", '/points/?itter=' + reqid , true );
            req.send();
        }
        catch (ex){
            error( ex );
        }

    });
    this.request++;

};