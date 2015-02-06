#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "poker.h"

int extract_rank(char c) {
  int rank;
  if (c == 't' || c == 'T')
    rank = 10;
  else if (c == 'j' || c == 'J')
    rank = 11;
  else if (c == 'q' || c == 'Q')
    rank = 12;
  else if (c == 'k' || c == 'K')
    rank = 13;
  else if (c == 'a' || c == 'A')
    rank = 14;
  else if (c >= '2' && c <= '9')
    rank = c - '0';
  else
    rank = 101;
  
  rank -= 2;
  return rank;
}

int extract_suit(char c) {
  int suit;
  if (c == 'c' || c == 'C')
    suit = CLUB;
  else if (c == 'd' || c == 'D')
    suit = DIAMOND;
  else if (c == 'h' || c == 'H')
    suit = HEART;
  else if (c == 's' || c == 'S')
    suit = SPADE;
  else
    suit = 99;
  return suit;
}

int check_format(char *str) {
  if (strlen(str) == 2)
    return 1; // true
  return 0; // false
}

void print_rank(int rank) {
  char *str;
  switch(rank) {
    case STRAIGHT_FLUSH:  str = "STRAIGHT_FLUSH";   break;
    case FOUR_OF_A_KIND:  str = "FOUR_OF_A_KIND";   break;
    case FULL_HOUSE:      str = "FULL_HOUSE";       break;
    case FLUSH:           str = "FLUSH";            break;
    case STRAIGHT:        str = "STRAIGHT";         break;
    case THREE_OF_A_KIND: str = "THREE_OF_A_KIND";  break;
    case TWO_PAIR:        str = "TWO_PAIR";         break;
    case ONE_PAIR:        str = "ONE_PAIR";         break;
    case HIGH_CARD:       str = "HIGH_CARD";        break;
    default:              str = "NOT_FOUND";        break;
  }
  printf("Rank: %s (%d)\n", str, rank);
}

int main(int argc, char *argv[]) {
  // printf("===== poker-eval =====\n");

  if (argc != 6) {
    printf("Usage: ./allfive card1 card2 card3 card4 card5\n");
    return -1;
  }

  if (!check_format(argv[1]) || !check_format(argv[2]) || !check_format(argv[3]) || !check_format(argv[4]) || !check_format(argv[5])) {
    printf("Invalid card format\n");
    return -1;
  }

  int deck[52];
  init_deck(deck);

  int r1 = extract_rank(argv[1][0]);
  int s1 = extract_suit(argv[1][1]);
  int r2 = extract_rank(argv[2][0]);
  int s2 = extract_suit(argv[2][1]);
  int r3 = extract_rank(argv[3][0]);
  int s3 = extract_suit(argv[3][1]);
  int r4 = extract_rank(argv[4][0]);
  int s4 = extract_suit(argv[4][1]);
  int r5 = extract_rank(argv[5][0]);
  int s5 = extract_suit(argv[5][1]);

  int c1_index = find_card(r1, s1, deck);
  int c2_index = find_card(r2, s2, deck);
  int c3_index = find_card(r3, s3, deck);
  int c4_index = find_card(r4, s4, deck);
  int c5_index = find_card(r5, s5, deck);

  if (c1_index == -1 || r1 == 99 || s1 == 99)
    printf("!!!!! c1_index is invalid !!!!!\n");
  if (c2_index == -1 || r2 == 99 || s2 == 99)
    printf("!!!!! c2_index is invalid !!!!!\n");
  if (c3_index == -1 || r3 == 99 || s3 == 99)
    printf("!!!!! c3_index is invalid !!!!!\n");
  if (c4_index == -1 || r4 == 99 || s4 == 99)
    printf("!!!!! c4_index is invalid !!!!!\n");
  if (c5_index == -1 || r5 == 99 || s5 == 99)
    printf("!!!!! c5_index is invalid !!!!!\n");

  int c1 = deck[c1_index];
  int c2 = deck[c2_index];
  int c3 = deck[c3_index];
  int c4 = deck[c4_index];
  int c5 = deck[c5_index];

  printf("Hand: %d-%d %d-%d %d-%d %d-%d %d-%d\n", r1, s1, r2, s2, r3, s3, r4, s4, r5, s5);

  int rank = hand_rank(eval_5cards(c1, c2, c3, c4, c5));

  print_rank(rank);

  return rank;
}