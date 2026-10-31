export function parseCompletionBlanks(question) {
	const title = question.title || ''
	const answerCount = question.answer_count || []
	const blanks = []
	const regex = /_+/g
	let match
	let blankIndex = 0
	while ((match = regex.exec(title)) !== null) {
		blanks.push({
			start: match.index,
			end: match.index + match[0].length,
			charCount: answerCount[blankIndex] || 1,
			index: blankIndex,
		})
		blankIndex++
	}
	return blanks
}

export function getCompletionParts(question, userAnswers) {
	const title = question.title || ''
	const blanks = parseCompletionBlanks(question)
	const parts = []
	let lastIndex = 0
	blanks.forEach((blank, blankIndex) => {
		if (blank.start > lastIndex) {
			parts.push({ type: 'text', text: title.substring(lastIndex, blank.start) })
		}
		const answers = userAnswers[question.id] || []
		const blankAnswers = answers[blankIndex] || []
		parts.push({
			type: 'blank',
			blankIndex,
			chars: blankAnswers.length > 0 ? [...blankAnswers] : new Array(blank.charCount).fill(''),
		})
		lastIndex = blank.end
	})
	if (lastIndex < title.length) {
		parts.push({ type: 'text', text: title.substring(lastIndex) })
	}
	return parts
}

export function isQuestionAnswered(q, answer) {
	if (q.type === 'completion') {
		if (!answer || answer.length === 0) return false
		return answer.every((blankAnswer) =>
			Array.isArray(blankAnswer) && blankAnswer.every((char) => char && char !== ''),
		)
	}
	if (q.type === 'single' || q.type === 'multiple') {
		return Array.isArray(answer) && answer.length > 0
	}
	if (q.type === 'text') {
		return !!(answer && String(answer).trim() !== '')
	}
	if (q.type === 'uploadImage' || q.type === 'uploadDocument') {
		return Array.isArray(answer) && answer.length > 0
	}
	return false
}

export function getUserAnswers(answer) {
	if (Array.isArray(answer)) {
		return answer.length > 0 ? answer[0] : null
	}
	return answer !== undefined && answer !== null ? answer : null
}

export function getUploadAnswers(answer) {
	return Array.isArray(answer) ? answer : []
}

export function getFileName(url) {
	if (!url) return ''
	const parts = url.split('/')
	const name = (parts[parts.length - 1] || '').split('?')[0]
	return name || '未知文件'
}
