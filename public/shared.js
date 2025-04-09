var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) =>
  key in obj
    ? __defProp(obj, key, {
        enumerable: true,
        configurable: true,
        writable: true,
        value,
      })
    : (obj[key] = value);
var __publicField = (obj, key, value) =>
  __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
const letters = ["A", "B", "C", "D"];

function generateQuestionElement(question, index, onSelect) {
  const div = document.createElement("div");
  div.className = "question border p-3 mb-3 rounded";
  div.innerHTML = "";
  if (question.question.paragraph != null) {
    div.innerHTML = `<p>${question.question.paragraph}</p>`;
  }

  if (question.topic == "math") {
    div.innerHTML += `<p><strong>${index + 1}) ${wrapMath(
      question.question.question
    )}</strong></p>`;
  } else {
    div.innerHTML += `<p><strong>${index + 1}) ${
      question.question.question
    }</strong></p>`;
  }

  for (let i = 0; i < question.question.choices.length; i++) {
    const optionWrapper = document.createElement("div");
    optionWrapper.className = "form-check";
    const input = document.createElement("input");
    input.className = "form-check-input d-none";
    input.type = "radio";
    input.id = `option-${question.id}-${i}`;
    input.name = `question-${question.id}`;
    input.value = i.toString();
    const label = document.createElement("label");
    label.className =
      "form-check-label btn btn-outline-primary w-100 text-start py-2";
    label.htmlFor = input.id;
    let choice = question.question.choices[i];
    if (question.topic == "math") {
      choice = wrapMath(choice);
    }
    label.innerHTML = `<strong>${letters[i]}</strong>: ${choice}`;
    if (question.user_answer === void 0 && onSelect) {
      input.addEventListener("change", () => {
        onSelect(i);
      });
    } else {
      input.checked = i == question.user_answer;
      label.classList.remove("btn-outline-primary", "btn-primary", "active");
      label.classList.add("btn-secondary", "disabled");
      if (input.checked && i != question.question.correct_answer) {
        label.classList.add("btn-danger");
      }
      if (i == question.question.correct_answer) {
        label.classList.add("btn-success");
      }
    }
    optionWrapper.appendChild(input);
    optionWrapper.appendChild(label);
    div.appendChild(optionWrapper);
  }
  if (question.user_answer > -1) {
    const rationale = document.createElement("p");
    if (question.topic == "math") {
      rationale.innerHTML = "<br>" + wrapMath(question.question.explanation);
    } else {
      rationale.innerHTML = "<br>" + question.question.explanation;
    }

    div.appendChild(rationale);
  }
  return div;
}

export function wrapMath(text) {
  // First process display math ($$...$$)
  let processed = text.replace(/\$\$([^$]+)\$\$/g, (_, math) => {
    try {
      return katex.renderToString(math, { displayMode: true });
    } catch (e) {
      console.error("Error rendering display math:", e);
      return `$$${math}$$`;
    }
  });

  // Then process inline math ($...$) but skip currency values
  processed = processed.replace(
    /(?<!\\)\$([^$\s]+(?:\s+[^$\s]+)*)\$(?!\d)/g,
    (_, math) => {
      // Skip if it looks like currency (space or punctuation after $)
      if (/^\d+\.?\d*\s*$/.test(math)) return `$${math}$`;

      try {
        return katex.renderToString(math, { displayMode: false });
      } catch (e) {
        console.error("Error rendering inline math:", e);
        return `$${math}$`;
      }
    }
  );

  // Process standalone LaTeX commands but skip when part of normal text
  processed = processed.replace(/(\\[a-zA-Z]+)(?![a-zA-Z])/g, (_, command) => {
    // Skip if part of a word (like in "problem_solving")
    if (/^\\[a-z]+_[a-z]+$/i.test(command)) return command;

    // Skip if in the middle of a word
    const prevChar = command[0] === "\\" ? "" : command[0];
    if (prevChar && /[a-zA-Z]/.test(prevChar)) return command;

    try {
      return katex.renderToString(command, { displayMode: false });
    } catch (e) {
      console.error("Error rendering LaTeX command:", e);
      return command;
    }
  });

  return processed;
}

class Score {
  constructor() {
    __publicField(this, "correct");
    __publicField(this, "total");
    this.correct = 0;
    this.total = 0;
  }
  getPercentage() {
    return this.total > 0 ? (this.correct / this.total) * 100 : 0;
  }
}
class AttemptResult {
  constructor() {
    __publicField(this, "topics");
    this.topics = {};
  }
  addTopic(name) {
    if (!(name in this.topics)) {
      this.topics[name] = new Score();
    }
  }
  getTotalScore() {
    let totalCorrect = 0;
    for (const score of Object.values(this.topics)) {
      totalCorrect += score.correct;
    }
    return totalCorrect;
  }
  getTotalScorePercentage() {
    let totalCorrect = 0;
    let totalQuestions = 0;
    for (const score of Object.values(this.topics)) {
      totalCorrect += score.correct;
      totalQuestions += score.total;
    }
    return totalQuestions > 0 ? (totalCorrect / totalQuestions) * 100 : 0;
  }
  static fromAnsweredQuestions(questions) {
    const result = new AttemptResult();
    for (const question of questions) {
      const topic = question.category;
      if (!topic) {
        console.warn(
          "invalid question. missing category",
          JSON.stringify(question)
        );
        continue;
      }
      result.addTopic(topic);
      result.topics[topic].total++;
      if (question.question.correct_answer == question.user_answer) {
        result.topics[topic].correct++;
      }
    }
    return result;
  }
}
function findQuestion(all_questions, question_id) {
  for (const category in all_questions) {
    if (Object.prototype.hasOwnProperty.call(all_questions, category)) {
      const questionsInCategory = all_questions[category];
      const found = questionsInCategory.find((q) => q.id === question_id);
      if (found) {
        return { ...found, category };
      }
    }
  }
  return null;
}
async function fetchQuestions() {
  const response = await fetch("/api/daily");
  if (!response.ok) {
    throw new Error(
      `Failed to fetch questions.json: ${response.status} ${response.statusText}`
    );
  }
  return response.json();
}
export {
  AttemptResult,
  Score,
  fetchQuestions,
  findQuestion,
  generateQuestionElement,
  letters,
};
