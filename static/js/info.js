

function Info(data, stream) {

    this.totalPts  = data.PointCnt || 1;
    this.loadedPts = 1;
    this.padding   = 10;
    this.progress  = document.createElement( 'canvas' );
    this.animation = -1;
    this.bpos      = 0;
    this.stroke_widths = {
        ready:  { width: 2,   color: '#0f0' },
        line:   { width: 1.5, color: '#000' },
        bar:    { width: 3,   color: '#00f' }
    };

    this.stroke_max = 0;
    for( s in this.stroke_widths ){
        this.stroke_max = Math.max( this.stroke_max, this.stroke_widths[s].width );
    }


    stream.on( 'pointbatch', this.progress_tick.bind( this ) );
    stream.on( 'done', this.end_progress.bind( this ) );

    let infoparts = document.getElementsByClassName( 'info' );
    ['No. pts', 'View center'].forEach(function(lbl, i){
        let span = document.createElement( 'span' );
        span.className = 'lbl';
        span.appendChild( document.createTextNode( lbl ) );
        infoparts[i].appendChild( span );
    });

    // left side total number of items
    let ncnt = document.createElement( 'div' );
    ncnt.appendChild( document.createTextNode( this.totalPts ) );
    infoparts[0].appendChild( ncnt );

    this.infodom = document.createElement( 'div' );
    this.infodom.id  = 'infoparts';
    ['x', 'y', 'z'].forEach(function(lbl){
        let wrap = document.createElement( 'div' );
        wrap.className = 'label-' + lbl;
        let dlbl = document.createElement( 'span' );
        dlbl.className = 'lbl';
        dlbl.appendChild( document.createTextNode( lbl ) );
        let dval = document.createElement( 'span' );
        dval.className = 'val';
        dval.appendChild( document.createTextNode( '0' ) );

        wrap.appendChild( dlbl );
        wrap.appendChild( dval );
        this.infodom.appendChild( wrap );
    }.bind( this ));
    infoparts[1].appendChild( this.infodom );

    // progress bar stuff, for now keep it simple
    let pwrapper = document.getElementById( 'progress' );
    if( pwrapper ){
        let rect = pwrapper.getBoundingClientRect();
        this.progress.width  = rect.width;
        this.progress.height = rect.height;
        pwrapper.appendChild( this.progress );
    }
    else {
        console.warn( 'where can progress canvas alternately go?' );
    }
};

Info.prototype.draw = function()
{
    this.bpos += 30;
    let ctx = this.progress.getContext( '2d' );
    let w   = this.progress.width;
    let y   = this.progress.height * .5;
    let p   = this.padding + 1;

    ctx.clearRect( p, y - this.stroke_max, w - 2 * p, 2 * this.stroke_max );
    
    // draw line, ready & bar
    let x = p;
    Object.keys( this.stroke_widths ).forEach( function(k){
        
        ctx.strokeStyle = this.stroke_widths[k].color;
        ctx.lineWidth = this.stroke_widths[k].width;
        ctx.beginPath();

        if( k == 'bar' ){
            let px = this.bpos % (w - p);
            ctx.moveTo( px+30, y - this.stroke_widths[k].width * .5 );
            ctx.lineTo( px-30, y - this.stroke_widths[k].width * .5 );
        }
        else {
            let len = 0;
            switch( k ){
                case 'line':
                    len = w - len - p;
                    break;
                case 'ready':
                    len = (this.loadedPts / this.totalPts) * (w - p);
                    break;
            }
            ctx.moveTo( x, y - this.stroke_widths[k].width * .5 );
            ctx.lineTo( len, y - this.stroke_widths[k].width * .5 );
            x += len;
        }
        ctx.stroke();
        ctx.moveTo( 0, 0 );
    }.bind( this ) );
    console.log( 'animation frame!' );
};

Info.prototype.start_progress = function()
{
    let ctx = this.progress.getContext( '2d' );
    let w   = this.progress.width;
    let h   = this.progress.height;
    let p   = this.padding;

    ctx.fillStyle = '#fff';
    ctx.fillRect( 0, 0, w, h );
    ctx.strokeStyle = '#000';
    ctx.strokeWidth = 1.2;

    ctx.save();
    [p, w-p].forEach( function(x){
        ctx.beginPath();
        ctx.moveTo( x, p );
        ctx.lineTo( x, h-p );
        ctx.stroke();
        ctx.restore();
    } );

    this.animation = setInterval( this.draw.bind( this ), 20 );

};

Info.prototype.end_progress = function()
{
    if( this.animation > 0 ){
        clearInterval( this.animation );
    }
    // a msg here?!
    // let ctx = this.progress.getContext( '2d' );
    // let w   = this.progress.width;
    // let h   = this.progress.height;
    // let p   = this.padding;

    // ctx.fillStyle = '#fff';
    // ctx.fillRect( 0, 0, w, h );

};

Info.prototype.progress_tick = function(cnt)
{
    this.loadedPts = cnt;
    // console.log( 'update progress tick', cnt );
};

Info.prototype.set_viewoport = function(viewport)
{
    Array.from( this.infodom.getElementsByTagName( 'div' ) )
         .forEach( function(cont, i){
            let valdoms = cont.getElementsByClassName( 'val' );
            if( valdoms.length == 1 ){
                valdoms[0].innerHTML = Math.trunc( viewport[i] * 10 ) / 10;
            }
         } );
};