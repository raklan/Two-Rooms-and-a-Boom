class Bomber {
    hasWon(player, gameState){
        let president = gameState.players.find(player => player.role === 'president')
        
        return player.room === president.room;
    }
}