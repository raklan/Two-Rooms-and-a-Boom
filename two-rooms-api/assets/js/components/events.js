class PlayerListUpdateEvent extends CustomEvent{    
    static EVENT_NAME = 'playerListUpdated'
    constructor(newPlayerList){
        super(
            PlayerListUpdateEvent.EVENT_NAME,
            {
                detail: { newPlayerList },
                bubbles: true,
                composed: true
            }
        );        
    }

    get newPlayerList(){
        return this.detail.newPlayerList;
    }
}

class PlayerSelectedEvent extends CustomEvent{
    static EVENT_NAME = 'playerSelected'
    constructor(selectedPlayer){
        super(
            PlayerSelectedEvent.EVENT_NAME,
            {
                detail: { selectedPlayer },
                bubbles: true,
                composed: true
            }
        )
    }

    get selectedPlayer(){
        return this.detail.selectedPlayer;
    }
}

export { PlayerListUpdateEvent, PlayerSelectedEvent }