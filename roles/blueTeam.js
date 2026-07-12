class BlueTeam {
    hasWon(player, gameState){
        let president = gameState.players.find(player => player.role === 'president')
        let bomber = gameState.players.find(player => player.role === 'bomber')
        
        return bomber.room !== president.room;
    }
}