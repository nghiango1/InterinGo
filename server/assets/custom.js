function addCommand() {
    document.getElementById('repl-output').append('>> ' + document.getElementById('repl-input').value + '\n');
}

function updateWrap() {
    document.getElementById('hiddenwrap').checked = document.getElementById('wrap').checked
}

function updateHide() {
    document.getElementById('hiddenhide').checked = document.getElementById('hide').checked
}

function copyCode(id) {
    document.getElementById('repl-input').value = document.getElementById(id).innerText
    document.getElementById('repl-send').click()
}

function scrollBottom() {
    let replOutput = document.getElementById('repl-output')
    replOutput.scrollTop = replOutput.scrollHeight
}
