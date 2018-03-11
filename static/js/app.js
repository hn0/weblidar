
var app = function() {
    var ccontainer = document.getElementById( 'cancontainer' );
    if( ccontainer ) {

        this.load_url( 'info' )
            .then(
                function(data){
                    console.log( 'got the data', data );
                    var stream = new Stream();
    	    	    var c      = new LidarCanvas( this, stream, ccontainer, data );
            }.bind( this ),
                function(){
                	console.warn( 'Application unable to start, cannot load the data' );
            });
    }
    else {
        console.warn( 'there should be an error response here!' );
    }
};

app.prototype.load_url = function(url, method, data) {

	if( !method ){
		method = "GET";
	}

	var p = new Promise(function(resolve,reject){

		var xhr = new XMLHttpRequest();
		xhr.open( method, url );
		xhr.onload = function(){
            resolve( xhr.responseText );
        };

        xhr.onerror = function() {
			reject();
		}

		xhr.send();
	});

	return p;
};