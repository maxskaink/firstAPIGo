import { Flex, Text, Spinner, Stack} from "@chakra-ui/react"
import { useQuery } from "@tanstack/react-query"
import TodoItem from "./TodoItem"
import { BASE_URL } from "../App"

export type Todo = {
    _id: number,
    body: string,
    completed: boolean
}

const TodoList = () => {
    const {data:todos, isLoading} = useQuery<Todo[]>({
        queryKey: ['todos'],
        queryFn: async () => {
            try {
                const res = await fetch(BASE_URL + "/todos");
                const data = await res.json();
                if (!res.ok) {
                    throw new Error(data.message)
                }
                return data || []
            } catch (error) {
                console.error(error)
            }
        }
    })

    return(
        <>
        <Text
            fontSize={"4x1"}
            textTransform={"uppercase"}
            fontWeight={"bold"}
            textAlign={"center"}
            my={2}
        >
            Todos task
        </Text>
        {isLoading && (
            <Flex justifyContent={"center"} my={4}>
                <Spinner size={"x1"}/>
            </Flex>
        )}
        {!isLoading && todos?.length === 0&&(
            <Stack alignItems={"center"} gap={3}>
                <Text 
                    fontSize={"x1"}
                    textAlign={"center"}
                    color={"gray.500"}
                > No todos found
                </Text>
            </Stack>
        )}
        <Stack gap={3}>
            {
                todos?.map(todo => (
                    <TodoItem key={todo._id} todo={todo} />
                ))
            }
        </Stack>
        </>
    )
}

export default TodoList