import pygame
import random
import math
import time

CLEAR_COLOR = (30,30,55)
MENU_STATE = 0
GAME_STATE = 1
OVER_STATE = 2

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
size = (1000,600)
screen = pygame.display.set_mode((size))
pygame.display.set_caption("Let's Tesselate! - A fictional tesselated multicore system")
pygame.mouse.set_visible(0)
game_font = pygame.font.Font("font.ttf", 35)
highscore_font = pygame.font.Font("font.ttf", 25)
clock = pygame.time.Clock()
highscore_file = "highscores"
screen.fill(CLEAR_COLOR)
rand = random.randrange

done = False
state = MENU_STATE
go_text = None
hs_text = None
highscore_text = None
highscores = {}
with open(highscore_file) as highscore_data:
	for line in highscore_data:
		score_difficulty = line.split()[0]
		score = line.split()[1]
		highscores[int(difficulty)] = int(score)
time_elapsed = 0

while not done:
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
			elif event.key == pygame.K_RETURN:
				if state == OVER_STATE:
					done = False
					state = GAME_STATE
					player1 = Player()
					bubbles = [Bubble()]
					boni = [Bonus()]
					points = 0
					go_text = None
					hs_text = None
					time_elapsed = 0
				elif state == MENU_STATE:
					state = GAME_STATE
					time_elapsed = 0

			else:
				if state == GAME_STATE:
					player1.change_rot()
				elif state == MENU_STATE:
					state = GAME_STATE
					time_elapsed = 0

	clock.tick(60)
	screen.fill(CLEAR_COLOR)
	if state == MENU_STATE:
		draw_menu()
		player1.draw()
		player1.move()
	elif state == GAME_STATE:
		player1.draw()
		if not player1.move():
			state = OVER_STATE
		do_bubbles()
		if check_collision(player1):
			state = OVER_STATE
	elif state == OVER_STATE:
		do_bubbles()
		if not go_text:
			go_text, hs_text = game_over()
		screen.blit(go_text, (size[0]/2-go_text.get_rect().width/2,size[1]/2-go_text.get_rect().height/2))
		screen.blit(hs_text, (size[0]/2-hs_text.get_rect().width/2,size[1]-size[1]/4-hs_text.get_rect().height/2))
	pygame.display.flip()
	time_elapsed+=clock.get_time()
