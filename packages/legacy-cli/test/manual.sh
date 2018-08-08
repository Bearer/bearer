set -e

export NC='\e[0m' # No Color
export BOLD='\e[1m' # No Color
export WHITE='\e[1;37m'
export BLUE='\e[0;34m'
export GREEN='\e[0;32m'
export DOWN='\e[B'

#################
# bearer --help #
#################
echo -e "${BLUE}${BOLD}TESTING:${NC} \`bearer --help\`"
bearer --help
echo -e "${GREEN}${BOLD}WORKS!!!${NC}"


##############
# bearer new #
##############
echo
echo -e "${BLUE}${BOLD}TESTING:${NC} \`bearer new TestingGoats\` OAuth2"
echo -ne '\n' | bearer new TestingGoats
echo -e "${GREEN}${BOLD}WORKS!!!${NC}"
rm -rf TestingGoats

echo
echo -e "${BLUE}${BOLD}TESTING:${NC} \`bearer new TestingGoats\` NoAuth"
echo -ne "${DOWN}\n" | bearer new TestingGoats
echo -e "${GREEN}${BOLD}WORKS!!!${NC}"
rm -rf TestingGoats

echo
echo -e "${BLUE}${BOLD}TESTING:${NC} \`bearer new TestingGoats\` apiKey"
echo -ne "${DOWN}${DOWN}\n" | bearer new TestingGoats
echo -e "${GREEN}${BOLD}WORKS!!!${NC}"
rm -rf TestingGoats

############
# bearer g #
############
echo
echo -e "${BLUE}${BOLD}TESTING:${NC} \`bearer g\` Views"
echo -ne '\n' | bearer new TestingGoats
cd TestingGoats
echo -ne "${DOWN}\rTest" | bearer g
echo -e "${GREEN}${BOLD}WORKS!!!${NC}"
cd .. && rm -rf TestingGoats
