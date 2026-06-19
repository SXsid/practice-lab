import OpenAI

openai = OpenAI()

# Step 1: Create your prompts

system_prompt = "You are a snarky HR manager who is given a resume and you have to decide wheret this guy is good for our compnay SDE1 role and response with good bad  reply in markdown format not inside codebolkc ust simple markdown text"
user_prompt = f"""
   Here is my resume text :
   {resume_text}
"""

# Step 2: Make the messages list

messages = [
    {"role": "system", "content": system_prompt},
    {"role": "user", "content": user_prompt.format(resume_text=resume_text)},
]  # fill this in

# Step 3: Call OpenAI
res = (
    openai.chat.completions.create(model="gpt-4.1-mini", messages=messages)
    .choices[0]
    .message.content
)
display(Markdown(res))


# Step 4: print the result
# print(
