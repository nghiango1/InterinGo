function addCommand() {
    document.getElementById('repl-output').append('>> ' + document.getElementById('repl-input').value + '\n');
}

function updateWrap() {
    document.getElementById('hiddenwrap').checked = document.getElementById('wrap').checked
}

function updateHide() {
    document.getElementById('hiddenhide').checked = document.getElementById('hide').checked
}

function linkedHide() {
    document.getElementById('linked').checked = !document.getElementById('linked').checked
}

function toggleStick() {
    document.getElementById('stick').checked = !document.getElementById('stick').checked 
}

function updateStick() {
    document.getElementById('hiddenstick').checked = document.getElementById('stick').checked
}

function copyCode(id) {
    document.getElementById('repl-input').value = document.getElementById(id).innerText
    document.getElementById('repl-send').click()
}

function clearEvalResult(id) {
    document.getElementById('repl-result').innerText = ""
}

function copyEvalResult(id) {
    document.getElementById('repl-output').append(document.getElementById('repl-result').innerText);
}

function scrollBottom() {
    let replOutput = document.getElementById('repl-output')
    replOutput.scrollTop = replOutput.scrollHeight
}
