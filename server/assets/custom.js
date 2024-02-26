function addCommand() {
    document.getElementById('repl-output').append('>> ' + document.getElementById('repl-input').value + '\n');
}

function updateWrap() {
    document.getElementById('hiddenwrap').checked = document.getElementById('wrap').checked
}

function scrollBottom() {
    let replOutput = document.getElementById('repl-output')
    replOutput.scrollTop = replOutput.scrollHeight
}
