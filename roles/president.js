class President {
    hasWon(player, gameState){
        let bomber = gameState.players.find(player => player.role === 'bomber')
        
        return player.room === bomber.room;
    }
}