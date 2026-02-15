import learn


def executer():

    result = learn.add_task.delay(2, 3)
    print(result.get(timeout=10))
    print(result.id)


if __name__ == "__main__":

    executer()
