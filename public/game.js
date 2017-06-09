$(document).ready(function(){
	//Canvas stuff
	var canvas = $("#canvas")[0];
	var ctx = canvas.getContext("2d");
	var w = $("#canvas").width();
	var h = $("#canvas").height();

	var socket = io();

	//Lets save the cell width in a variable for easy control
	var cw = 10;
	var food;
	var score = 0;

	// Multiple sankes
	var number_of_snakes = 0;
	
	//Lets create the snake now
	var snakes = []; // array of snakes
	// var snake_array; //an array of cells to make up the snake
	
	function init()
	{
		snakes = [];
		// create_snake();
		// create_food(); //Now we can see the food particle
		//finally lets display the score
		score = 0;
		
		//Lets move the snake now using a timer which will trigger the paint function
		//every 60ms
		/*if(typeof game_loop != "undefined") clearInterval(game_loop);
		game_loop = setInterval(paint, 60);*/
	}
	//init();


	socket.on('new_state', function(msg) {
		console.log('Received state: ' + msg)
		// Parses msg content to snakes
		let message = JSON.parse(msg);
		console.log(message);
		snakes = JSON.parse(message.snakes);
		food = JSON.parse(message['food']);
		console.log(snakes);
		paint();
	});

	console.log('Socket shit completed')
	
	//Lets create the food now
	function create_food()
	{
		food = {
			x: Math.round(Math.random()*(w-cw)/cw), 
			y: Math.round(Math.random()*(h-cw)/cw), 
		};
		//This will create a cell with x/y between 0-44
		//Because there are 45(450/10) positions accross the rows and columns
	}
	
	//Lets paint the snake now
	function paint()
	{
		//To avoid the snake trail we need to paint the BG on every frame
		//Lets paint the canvas now
		ctx.fillStyle = "white";
		ctx.fillRect(0, 0, w, h);
		ctx.strokeStyle = "black";
		ctx.strokeRect(0, 0, w, h);
		
		food_loc = food.location
		// console.log('Food location: ' + food_loc)
		paint_cell(food_loc.x, food_loc.y, "LightCoral")

		//The movement code for the snake to come here.
		//The logic is simple
		//Pop out the tail cell and place it infront of the head cell
		for (let i = 0; i < snakes.length; i++) {
			let s = snakes[i];
			// let nx = s.snake_array[0].x;
			// let ny = s.snake_array[0].y;
		
		
			//These were the position of the head cell.
			//We will increment it to get the new head position
			//Lets add proper direction based movement now
			/*
			if (s.d == "right") nx++;
			else if(s.d == "left") nx--;
			else if(s.d == "up") ny--;
			else if(s.d == "down") ny++;
			*/
		
			//Lets add the game over clauses now
			//This will restart the game if the snake hits the wall
			//Lets add the code for body collision
			//Now if the head of the snake bumps into its body, the game will restart
			/*
			if(nx == -1 || nx == w/cw || ny == -1 || ny == h/cw || check_collision(nx, ny, s.snake_array))
			{
				//restart game
				init();
				//Lets organize the code a bit now.
				return;
			}
			*/
		
			//Lets write the code to make the snake eat the food
			//The logic is simple
			//If the new head position matches with that of the food,
			//Create a new head instead of moving the tail
			/*
			if (nx == food.x && ny == food.y)
			{
				var tail = {x: nx, y: ny};
				score++;
				//Create new food
				create_food();
			}
			else
			{
				var tail = s.snake_array.pop(); //pops out the last cell
				tail.x = nx; tail.y = ny;
			}
			*/
			//The snake can now eat the food.
			
		
			// s.snake_array.unshift(tail); //puts back the tail as the first cell
			for(let j = 0; j < s.body.length; j++)
			{
				let c = s.body[j];
				//Lets paint 10px wide cells
				paint_cell(c.x, c.y, "LightSkyBlue");
			}
		
			//Lets paint the score
			var score_text = "Score: " + score;
			ctx.fillText(score_text, 5, h-5);
		}
	}
	
	//Lets first create a generic function to paint cells
	function paint_cell(x, y, color)
	{
		ctx.fillStyle = color;
		ctx.fillRect(x*cw, y*cw, cw, cw);
		ctx.strokeStyle = "white";
		ctx.strokeRect(x*cw, y*cw, cw, cw);
	}
	
	/*function check_collision(x, y, array)
	{
		//This function will check if the provided x/y coordinates exist
		//in an array of cells or not
		for(var i = 0; i < array.length; i++)
		{
			if(array[i].x == x && array[i].y == y) {
				return true;
			}
		}
		return false;
	}*/
	
	//Lets add the keyboard controls now
	$(document).keydown(function(e){
		var key = e.which;
		// TODO: convert this so that it emits the direction

		//We will add another clause to prevent reverse gear
		/*if(key == "37" && snakes[0].d != "right") snakes[0].d = "left";
		else if(key == "38" && snakes[0].d != "down") snakes[0].d = "up";
		else if(key == "39" && snakes[0].d != "left") snakes[0].d = "right";
		else if(key == "40" && snakes[0].d != "up") snakes[0].d = "down";*/
		//The snake is now keyboard controllable
		let leftKey = 37
		let upKey = 38
		let rightKey = 39
		let downKey = 40

		switch (key) {
			case leftKey: // d
				socket.emit('move', 'left')
				break;
			case rightKey: // s
				socket.emit('move', 'right') 
				break;
			case upKey: // a 
				socket.emit('move', 'up')
				break;
			case downKey: // w
				socket.emit('move', 'down')
				break;
		}
	})
})