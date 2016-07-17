var game;
function preload(){
  game = loadStrings("game.txt");
}
var dim = [];
var columns;
var rows;

function setup(){
  dim = split(game[0], " ");
  columns = int(dim[0]);
  rows = int(dim[1]);
  createCanvas(window.innerWidth, window.innerHeight);
}

function draw(){
  background(50);
  text(dim,20,20);
  buffer = 50;
  //board height & width
  var bh = height-2*buffer
  var bw = width-2*buffer

  if(rows > columns){
    var size = bh/rows;
  } else {
    var size = bw/columns;
  }
  push();
  translate(width/2-(size*columns)/2, height/2-(size*rows)/2);
  for(var i = 0; i < columns; i++){
    for(var j = 0; j < rows; j++){
      rect(i*size,j*size,size,size);
    }
  }
  pop();
}
