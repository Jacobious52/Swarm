var game;
function preload(){
  game = loadStrings("game.txt");
}
var dim = [];
var columns;
var rows;
var currentMove;
var moveCount;
var winner;
var gameRunning;

function setup(){
  if(game[game.length-1] == "B"){
    winner = "Blue"
  } else if (game[game.length-1] == "Y"){
    winner = "Yellow"
  } else {
    winner = "You all lose!"
  }
  gameRunning = true;
  currentMove = 1;
  moveCount = (game.length)-1;
  dim = split(game[0], " ");
  columns = int(dim[0]);
  rows = int(dim[1]);
  createCanvas(window.innerWidth, window.innerHeight);
}

var cTime;
var pTime;

function draw(){
  background(50);
  fill(255);
  text(dim[0] + " x " + dim[1],20,20);
  text("total number of moves: " + moveCount,20 ,40);
  text("move: " + currentMove,20 ,60);
  if(gameRunning){
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
        switch(game[currentMove].charAt(j*rows+i)){
          case "Y":
            fill(255,255,100);
          break;
          case "B":
            fill(100,175,255);
          break;
          case "0":
            fill(0,0);
          break;
          case "X":
            fill(100);
        }
        stroke(250);
        rect(i*size,j*size,size,size);
      }
    }
    pop();
  } else {
    push();
    textAlign(CENTER, CENTER);
    translate(width/2, height/2);
    textSize(30);
    text("THE WINNER IS:",0 ,0);
    textSize(50);
    text(winner,0 ,60);
    pop();
  }
  cTime = second();
  if(currentMove > moveCount-1){
    gameRunning = false;
  } else {
    //if a second has passed
    if(pTime != cTime){
      currentMove++;
    }
    pTime = cTime;
  }
}
