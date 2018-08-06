import pygame
import random
import math
import time

MENU_STATE = 0
GAME_STATE = 1
OVER_STATE = 2

def hex_to_rgb(hex):
	return tuple(int(hex.lstrip('#')[i:i+2], 16) for i in (0, 2 ,4))

def num_to_cmd(num):
	if num == 0:
		return "NOP"
	elif num == 1:
		return "NEG"
	elif num == 2:
		return "ADD"
	elif num == 3:
		return "SUB"
	elif num == 4:
		return "MUL"
	elif num == 5:
		return "DIV"
	elif num == 6:
		return "MOV"
	elif num == 7:
		return "JGZ"
	elif num == 8:
		return "JEZ"
	elif num == 9:
		return "OUT"
	elif num == 10:
		return "AND"
	elif num == 11:
		return "IOR"
	elif num == 12:
		return "XOR"
	else:
		return "???"

def dirsymb_to_arrowpos(symb, tile_size_x, tile_size_y):
	dir = None
	for i in symb:
		if i=="^":
			dir = (tile_size_x/2, 5)
		elif i==">":
			dir = (tile_size_x-5, tile_size_y/2)
		elif i=="v":
			dir = (tile_size_x/2, tile_size_y-5)
		elif i=="<":
			dir = (5, tile_size_y/2)
		else:
			continue
		break
	return dir

def to_time(time):
	minutes = str(int(time/1000/60))
	if len(minutes)==1:
		minutes = "0"+minutes
	seconds = str(int((time/1000)%60))
	if len(seconds)==1:
		seconds = "0"+seconds
	milliseconds = str(int((time/10)%100))
	if len(milliseconds)==1:
		milliseconds = "0"+milliseconds
	return str(minutes)+":"+str(seconds)+":"+str(milliseconds)

pygame.init()
size = (920,920)
screen = pygame.display.set_mode((size))
pygame.display.set_caption("Let's Tesselate!")
pygame.mouse.set_visible(0)
big_font = pygame.font.Font("font.ttf", 35)
small_font = pygame.font.Font("font.ttf", 25)
clock = pygame.time.Clock()
rand = random.randrange

line=""
while line!="RENDER_INFO_BEGIN":
	line = input()

width, height, inactive_color = [x for x in input().split()]
width, height = int(width), int(height)
tile_size_x = size[0]/width
tile_size_y = size[1]/height
inactive_color = hex_to_rgb(inactive_color)

screen.fill(inactive_color)

cmd_render_dict = {}
for i in range(13):
	cmd_render = small_font.render(num_to_cmd(i), 1, (255,255,255))
	mid_x, mid_y = (i/2 for i in cmd_render.get_size())
	cmd_render_dict[num_to_cmd(i)] = (cmd_render, (mid_x, mid_y))

single_steps = True
next_step = False
done = False
time_elapsed = 0
while not done:
	#keystrokes
	for event in pygame.event.get():
		if event.type == pygame.QUIT:
			done = True
		elif event.type == pygame.KEYDOWN:

			if event.key == pygame.K_ESCAPE:
				done = True
			elif event.key == pygame.K_F12:
				filename=str(rand(10000000))+'.png'
				print("saved as ", (filename))
				pygame.image.save(screen,filename)
			elif event.key == pygame.K_p:
				single_steps = not single_steps
				next_step = False
			elif event.key == pygame.K_RETURN:
				if single_steps:
					next_step = True
	#pygame stuff
	clock.tick(60)
	screen.fill(inactive_color)
	#graphics 1 - background for inactives

	#inputtext
	if (not single_steps) or (single_steps and next_step):
		line = input()
		while line != "###":
			if not line or line == "RENDER_INFO_END":
				done = True
				break
			line = line.split()
			if line[0] == "OUT:":
				pass
			else:
				core_num = int(line[0])
				core_x = int(core_num%width)*tile_size_x
				core_y = int(core_num/height)*tile_size_y
				core_col = line[1]
				core_col = hex_to_rgb(core_col)
				command = int(line[2][1])
				command = num_to_cmd(command)
				command_dir = line[5]
				arrowpos = dirsymb_to_arrowpos(command_dir, tile_size_x, tile_size_y)
				# print(command_dir)
				pygame.draw.rect(screen, core_col, (core_x,core_y,tile_size_x,tile_size_y))
				text_x = core_x+tile_size_x/2-cmd_render_dict[command][1][0]
				text_y = core_y+tile_size_y/2-cmd_render_dict[command][1][1]
				if arrowpos:
					pygame.draw.circle(screen, tuple((x-50)%255 for x in core_col), (int(core_x+arrowpos[0]),int(core_y+arrowpos[1])),7)
				screen.blit(cmd_render_dict[command][0], (text_x,text_y))

			line = input()
		#graphics 3 - grid
		for x in range(width):
			pygame.draw.line(screen, (0,0,0), (x*tile_size_x,0), (x*tile_size_x,size[1]))
		for y in range(height):
			pygame.draw.line(screen, (0,0,0), (0,y*tile_size_y), (size[0],y*tile_size_y))
		#pygame stuff again
		pygame.display.flip()
		if next_step:
			next_step = False

	time_elapsed+=clock.get_time()
