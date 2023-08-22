// https://sigaa.unb.br/sigaa/public/turmas/listar.jsf
// extrai as turmas de matérias ofertadas.
var lines = document.querySelectorAll(
  '#turmasAbertas tbody tr'
);

var info = new Map();
var current = '';

var res = [];

for (const element of lines) {
  if (element.classList.contains('agrupador')) {
    current = element.innerText.trim();
  } else {
    var s = element.innerText.split('\t').map(x => x.trim())
    res.push({ materia: current, turma: s[0], professor: s[2], horario: s[3] })
  }
}

res;
