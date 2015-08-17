/*
 * main.c
 *
 * Copyright (c) 2014 Jeremy Garff <jer @ jers.net>
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted
 * provided that the following conditions are met:
 *
 *     1.  Redistributions of source code must retain the above copyright notice, this list of
 *         conditions and the following disclaimer.
 *     2.  Redistributions in binary form must reproduce the above copyright notice, this list
 *         of conditions and the following disclaimer in the documentation and/or other materials
 *         provided with the distribution.
 *     3.  Neither the name of the owner nor the names of its contributors may be used to endorse
 *         or promote products derived from this software without specific prior written permission.
 * 
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
 * FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT
 * OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */


#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <signal.h>

#include "clk.h"
#include "gpio.h"
#include "dma.h"
#include "pwm.h"

#include "ws2811.h"

#include <sys/socket.h>
#include <sys/un.h>

#define ARRAY_SIZE(stuff)                        (sizeof(stuff) / sizeof(stuff[0]))

#define TARGET_FREQ                              WS2811_TARGET_FREQ
#define GPIO_PIN                                 18
#define DMA                                      5

#define WIDTH                                    336
#define HEIGHT                                   1
#define LED_COUNT                                (WIDTH * HEIGHT)

void error(const char *);

ws2811_t ledstring =
{
    .freq = TARGET_FREQ,
    .dmanum = DMA,
    .channel =
    {
        [0] =
        {
            .gpionum = GPIO_PIN,
            .count = LED_COUNT,
            .invert = 0,
            .brightness = 255,
        },
        [1] =
        {
            .gpionum = 0,
            .count = 0,
            .invert = 0,
            .brightness = 0,
        },
    },
};

ws2811_led_t matrix[WIDTH][HEIGHT];

void matrix_render(void)
{
    int x, y;

    for (x = 0; x < WIDTH; x++)
    {
        for (y = 0; y < HEIGHT; y++)
        {
            ledstring.channel[0].leds[(y * WIDTH) + x] = matrix[x][y];
        }
    }
}
/*
void matrix_raise(void)
{
    int x, y;

    for (y = 0; y < (HEIGHT - 1); y++)
    {
        for (x = 0; x < WIDTH; x++)
        {
            matrix[x][y] = matrix[x][y + 1];
        }
    }
}
*/
void parseRGBvalues(ws2811_led_t* leds, char* rgbValues, int nrOfLeds) {
	//printf("buf: %2x%2x%2x ", rgbValues[0], rgbValues[1], rgbValues[2]);
        int i;
	for(i = 0; i < nrOfLeds; i++){
			leds[i] = rgbValues[i*3]<<16 | rgbValues[(i*3)+1]<<8 | rgbValues[(i*3)+2];
			//leds[i] = 0x11<<16 | 0x11<<8 | 0x11;
	}
	return;
}

void getMatrixGOing(ws2811_led_t* leds){
    int x, y;
    for (x = 0; x < WIDTH; x++){
        for (y = 0; y < HEIGHT; y++){
            matrix[x][y] = leds[x];
        }
    }
    matrix_render();
    return;
}
/*
int dotspos[] = { 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13};
ws2811_led_t dotcolors[] =
{
    0x111111,
    0xFF0000,
    0x200000,  // red
    0x201000,  // orange
    0x202000,  // yellow
    0x002000,  // green
    0x00FF00,
    0x002020,  // lightblue
    0x000020,  // blue
    0x0000FF,
    0x100010,  // purple
    0x200010,  // pink
    0xFFFFFF,
    0xFFFFFF, 
};

void matrix_bottom(void)
{
    int i;

    for (i = 0; i < ARRAY_SIZE(dotspos); i++)
    {
        dotspos[i]++;
        if (dotspos[i] > (WIDTH - 1))
        {
            dotspos[i] = 0;
        }

        matrix[dotspos[i]][HEIGHT - 1] = dotcolors[i];
    }
}
*/
static void ctrl_c_handler(int signum){
    ws2811_fini(&ledstring);
}

static void setup_handlers(void){
    struct sigaction sa =
    {
        .sa_handler = ctrl_c_handler,
    };

    sigaction(SIGKILL, &sa, NULL);
}

void error(const char *msg){
    perror(msg);
    exit(0);
}

int main(int argc, char *argv[])
{
	printf("Good Day!\n");
    	int ret = 0;
    	//Nr of total leds on the cube
    	ws2811_led_t leds[336];
	setup_handlers();
    	printf("After Handler setup...\n");
    	if (ws2811_init(&ledstring))
        	return -1;
    	printf("Socket stuff...\n");
//--------------------------------------------------------
 
	int sockfd, newsockfd, servlen; //n for read bytes from socket
   	socklen_t clilen;
   	struct sockaddr_un  cli_addr, serv_addr;
   	//1008 = Nr of LEDs * 3 byte
   	char buf[1008];

   	if ((sockfd = socket(AF_UNIX,SOCK_STREAM,0)) < 0)
       		error("creating socket");
   	bzero((char *) &serv_addr, sizeof(serv_addr));
   	serv_addr.sun_family = AF_UNIX;
   	strcpy(serv_addr.sun_path, "/tmp/so");
   	servlen=strlen(serv_addr.sun_path) + sizeof(serv_addr.sun_family);
   	if(bind(sockfd,(struct sockaddr *)&serv_addr,servlen)<0)
       		error("binding socket"); 

	listen(sockfd,5);
	clilen = sizeof(cli_addr);
	newsockfd = accept(
        sockfd,(struct sockaddr *)&cli_addr,&clilen);
   	if (newsockfd < 0)
        	error("accepting");
//--------------------------------------------------------
	printf("Entering while loop...\n");
	int loop = 0;
    	while (1){
		loop++;
//       	printf("Inside while...\n");
//       	matrix_raise();
//       	printf("After raise...\n");
//       	matrix_bottom();
//       	printf("After bottom...\n");
//       	printf("After render...\n");

//----------------------------------------------------
		//printf("Before Send...\n");
 		send(newsockfd,"x",1,0);
		//printf("Before Receive...\n");
   		//see buffersize for magic nr 1008
      		recv(newsockfd,buf,1008,0); //n=     returns Nr read bytes
   		parseRGBvalues(leds, buf, 336);
   		getMatrixGOing(leds);
   		//matrix_render();

		/*int p;
  		for(p = 0; p < 1008; p++){
  			printf("%x ", buf[p]);
  		}*/
		//if(loop%1000 == 0)
  		//printf(" loop: %d\n",loop);

   		//printf("First buffer Values are: 1: %6x 2: %6x 3: %6x\n", buf[0], buf[1], buf[2]);
   		//printf("-\n");
//-----------------------------------------------------
        	if (ws2811_render(&ledstring)){
            		ret = -1;
            		break;
        	}

        	// 15 frames /sec
        	//usleep(1000000 / 21);
    	}
    	ws2811_fini(&ledstring);
	//close socket connection
    	close(newsockfd);
    	close(sockfd);
    	return ret;
}

