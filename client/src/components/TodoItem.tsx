 import { Badge, Box, Flex, Spinner, Text } from "@chakra-ui/react";
 import { FaCheckCircle } from "react-icons/fa";
 import { MdDelete } from "react-icons/md";
import { Todo } from "./TodoList";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../App";

 const TodoItem = ({ todo }: { todo: Todo }) => {
    const queryClient = useQueryClient()

	const {mutate: updateTodo, isPending:isUpdating} = useMutation({
		mutationKey: ["updateTodo"],
		mutationFn: async () => {
			try {
				const res =await  fetch(BASE_URL + "/todos/" + todo._id, {
					method: "PATCH",
					headers: {
						"Content-Type": "application/json"
					}
				})
				const data = await res.json()
				if (!res.ok) {
					throw new Error(data.message)
				}
				return data 
			} catch (error) {
				console.error(error)
			}
		},
        onSuccess: ()=> {
            queryClient.invalidateQueries({queryKey:["todos"]})
        }
	})
	
 	return (
 		<Flex gap={2} alignItems={"center"}>
 			<Flex
 				flex={1}
 				alignItems={"center"}
 				border={"1px"}
 				borderColor={"gray.600"}
 				p={2}
 				borderRadius={"lg"}
 				justifyContent={"space-between"}
 			>
 				<Text
 					color={todo.completed ? "green.200" : "yellow.100"}
 					textDecoration={todo.completed ? "line-through" : "none"}
 				>
 					{todo.body}
 				</Text>
 				{todo.completed && (
 					<Badge ml='1' colorScheme='green'>
 						Done
 					</Badge>
 				)}
 				{!todo.completed && (
 					<Badge ml='1' colorScheme='yellow'>
 						In Progress
 					</Badge>
 				)}
 			</Flex>
 			<Flex gap={2} alignItems={"center"}>
 				<Box color={"green.500"} cursor={"pointer"} onClick={() => updateTodo()}>
 					{!isUpdating && <FaCheckCircle size={25} />}
					{isUpdating && <Spinner size={"sm"} />}
 				</Box>
 				<Box color={"red.500"} cursor={"pointer"}>
 					<MdDelete size={25} />
 				</Box>
 			</Flex>
 		</Flex>
 	);
 };


export default TodoItem