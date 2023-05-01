const server = "https://sdscoep.club/engine/api/"
const focus = { row: -1, col: -1 }
let gameId = new Date().getTime()
let toImgCache = ""
let board = [["R'", "N'", "B'", "K'", "Q'", "B'", "N'", "R'"],
["P'", "P'", "P'", "P'", "P'", "P'", "P'", "P'"],
["-", "-", "-", "-", "-", "-", "-", "-"],
["-", "-", "-", "-", "-", "-", "-", "-"],
["-", "-", "-", "-", "-", "-", "-", "-"],
["-", "-", "-", "-", "-", "-", "-", "-"],
["P", "P", "P", "P", "P", "P", "P", "P"],
["R", "N", "B", "K", "Q", "B", "N", "R"]]

const paintBoard = board => {
    Array.from(document.getElementsByClassName("piece")).forEach(ele => ele.remove());
    for (let i = 0; i < 8; i++) {
        for (let j = 0; j < 8; j++) {
            const squareImg = document.getElementById(`${i}${j}`).innerHTML
            const pieceImg = getPieceImage(getPieceFromString(board[i][j]))
            document.getElementById(`${i}${j}`).innerHTML = pieceImg + squareImg
        }
    }
}

const getPieceImage = img => {
    if (!img) {
        return ""
    }
    return `
    <img class="piece" src="${img}" alt="${img.split("_")[1]}" />
    `
}

const getSquareImage = img => {
    if (!img) {
        return ""
    }
    return `
    <img class="square" src="${img}" alt="${img.split(" ")[2]} square" />
    `
}

const getPieceFromString = piece => {
    switch (piece) {
        case "K'":
            return "1x/b_king_1x.png"
        case "Q'":
            return "1x/b_queen_1x.png"
        case "R'":
            return "1x/b_rook_1x.png"
        case "B'":
            return "1x/b_bishop_1x.png"
        case "N'":
            return "1x/b_knight_1x.png"
        case "P'":
            return "1x/b_pawn_1x.png"
        case "K":
            return "1x/w_king_1x.png"
        case "Q":
            return "1x/w_queen_1x.png"
        case "R":
            return "1x/w_rook_1x.png"
        case "B":
            return "1x/w_bishop_1x.png"
        case "N":
            return "1x/w_knight_1x.png"
        case "P":
            return "1x/w_pawn_1x.png"
        default:
            return ""
    }
}

document.getElementById("chess-container").addEventListener("click", async (event) => {
    const isImg = event.target.nodeName === 'IMG';
    if (!isImg) {
        return;
    }

    const id = event.target.parentElement.id
    const row = Number(id[0])
    const col = Number(id[1])

    console.log(board[row][col]);

    if (focus.row != -1 && (board[row][col] == '-' || board[row][col].length == 2)) {
        from = { row: focus.row, col: focus.col }
        to = { row: row, col: col }
        makeMove(from, to)
        clearFocus()
        await play(from, to)
    }
    else if (board[row][col] != '-' && board[row][col].length != 2) {
        clearFocus()
        focus.row = row
        focus.col = col
        applyFocus()
    }
    else {
        clearFocus()
    }

});

const play = async (from, to) => {
    document.getElementById("loader").style.visibility = "visible";
    console.log(JSON.stringify({ "fromRow": from.row, "fromCol": from.col, "toRow": to.col, "toCol": to.col }))
    await fetch(server + "?id=" + gameId,
        {
            method: 'POST',
            headers: { 'Content-Type': 'application/json; charset=utf-8' },
            body: JSON.stringify({ "FromRow": from.row, "FromCol": from.col, "ToRow": to.row, "ToCol": to.col }),
        }).then(async resp => {
            switch (resp.status) {
                case 406: //not acceptable
                    makeMove(to, from,true)
                    alert(await resp.text())
                    break
                case 206: //non authoritative info
                    makeMove(to, from,true)	
                    alert(await resp.text())
                    break
                case 403: /// forbidden
                    makeMove(to, from,true)
                    alert(await resp.text())
                    break
                case 204: //no content	
                    alert(await resp.text())
                    break
                case 400:
                    alert("Oops something went wrong in the request")
                    break
                case 200: //OK
                    const data = await resp.json()
                    board = data['Board']
                    paintBoard(data['Board'])
		    setTimeout(()=>{
				if(data['Check']){
				alert("CHECK !")
			    }
			},500);

		    setTimeout(()=>{
				if(data['Mate']){
	 			alert("And Mate :| ");
			    }
			},500);
		    
                    break
                default:
                    alert("Oops ! Cannot handle unexpected response")
            }
        }).catch(err => {
            console.log(err);
        });
    document.getElementById("loader").style.visibility = "hidden";

}

const makeMove = (from, to, revert = false) => {
    console.log(from, to)
    const fromSquare = document.getElementById(`${from.row}${from.col}`)
    const toSquare = document.getElementById(`${to.row}${to.col}`)

    const fromImg = fromSquare.getElementsByClassName("piece")[0].src
    const fromBackGround = fromSquare.getElementsByClassName("square")[0].src
    const toBackGround = toSquare.getElementsByClassName("square")[0].src
	
    let cache = toSquare.getElementsByClassName("piece")[0]?.src

    document.getElementById(`${to.row}${to.col}`).innerHTML = getPieceImage(fromImg) + getSquareImage(toBackGround)
    document.getElementById(`${from.row}${from.col}`).innerHTML = ( revert && toImgCache ? getPieceImage(toImgCache): "" ) + getSquareImage(fromBackGround)
    toImgCache = cache 
}

const applyFocus = () => {
    const square = document.getElementById(`${focus.row}${focus.col}`)
    let background = square.getElementsByClassName("square")[0].src
    square.getElementsByClassName("square")[0].src = background.replace("gray", "brown")
}

const clearFocus = () => {
    if (focus.row == -1) return;
    const square = document.getElementById(`${focus.row}${focus.col}`)
    let background = square.getElementsByClassName("square")[0].src
    square.getElementsByClassName("square")[0].src = background.replace("brown", "gray")
    focus.row = -1
    focus.col = -1
}

const fetchBoard = async () => {
    document.getElementById("loader").style.visibility = "visible";

    await fetch(server + "?id=" + gameId,
        {
            method: 'GET',

        }).then(async resp => {
            switch (resp.status) {

                case 200: //OK
                    const data = await resp.json()
                    board = data['Board']
                    paintBoard(data['Board'])
                    return
                default:
                    alert("Cannot handle unexpected response")
            }
        }).catch(err => {
            console.log(err);
        })
    document.getElementById("loader").style.visibility = "hidden";

}



const initialise = async () => {
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });
    if (params.id) {
        gameId = params.id
    } else {
        gameId = new Date().getTime()
        window.location = window.location + "?id=" + gameId
    }
    await fetchBoard()
}

(async () => initialise())()
