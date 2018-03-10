
function Stream (){
    this.request = 0;
    this.points  = [];

    this.get()
        .then(function(data) {

            if( data.byteLength ){

                // first byte always is length of response
                var len = new Int32Array( data.slice( 0, 4 ) );

                console.log( 'got the array!!!', data, len, len[0] )
                // var dw = new DataView( data.slice( i, i + buf[1] ) );

                // if( data == -1 ){
                //     console.log('done and ready');
                // }
                // this.get() if returned array len > 1
            }

        }, function(){
            console.log( 'we got an error!' );
        });

}


Stream.prototype.get = function(){

    // we will need a response size!?
    this.request++;
    return new Promise( (success, error) => {
            
        var req = new XMLHttpRequest();
        req.responseType = 'arraybuffer';
        req.onload  = () => success( req.response );
        req.onerror = () => error( null );

        try{
            // maybe some kind of clientid will be needed!
            req.open( "GET", '/points/' + this.request , true );
            req.send();
        }
        catch (ex){
            error( ex );
        }

    });

}