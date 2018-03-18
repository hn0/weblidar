

function Info(data, stream) {

    this.totalPts  = data.PointCnt || 1;
    this.loadedPts = 1;
    this.padding   = 10;
    this.progress  = document.createElement( 'canvas' );
    this.animation = -1;
    this.bpos      = 0;
    this.stroke_widths = {
        line:   { width: 1.5, color: '#000' },
        ready:  { width: 2,   color: '#0f0' },
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
    let ctx = this.progress.getContext( '2d' );
    let w   = this.progress.width;
    let y   = this.progress.height * .5;
    let p   = this.padding + 1;


    ctx.fillStyle = '#ff0';
    ctx.fillRect( p, y - this.stroke_max * .5, w - 2 * p, this.stroke_max );

    // draw line, ready & bar

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
    // console.log( 'update progress tick', cnt );
};

Info.prototype.set_viewoport = function(viewport)
{

};