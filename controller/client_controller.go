package controller

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/hirepb"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/model"
	"google.golang.org/grpc"
)

var ServiceClient hirepb.HireDataServiceClient

var mainMenu = `
=================================
Hiring DataBase
=================================

Main Menu:
1) Show Hires
2) Create Hire
3) Find Hire
4) Update Hired Info
5) Delete Hire
6) Quit

7) Add Dummy Employees

Enter option [1-7]: `

var findMenu = `
Find hires by:
1) Name
2) Type
3) Role
4) Duration
5) Tags
6) Back

Enter option [1-6]: `

var (
	server_conn        *grpc.ClientConn
	createHireName     = "Enter name: "
	createHireType     = "Enter type [Employee or Contractor]: "
	createHireDuration = "Enter contract duration: "
	createHireRole     = "Enter employee role: "
	createHireTags     = "Enter tags [Separate with commas]: "
	findTxt            = "Enter text to find: "
	deleteTxt          = "Enter the name of the hire to delete: "
	confirmDeleteTxt   = "Confirm delete of hire? [y/n]:"
	updateTxt          = "Enter the name of the hire to update: "
)

// DisplayMainMenu does the initial view for the client
func DisplayMainMenu() {
	var user_input int
	for {
		fmt.Printf("%s", mainMenu)
		fmt.Scanf("%d", &user_input)
		if validateMainMenuInput(user_input) {
			switch user_input {
			case 1:
				listHiresMenu()
			case 2:
				createHireMenu()
			case 3:
				findHiresMenu()
			case 4:
				updateHireData()
			case 5:
				deleteHire()
			case 6:
				return
			case 7:
				addDummyEmployees()
			}

		} else {
			fmt.Println("Invalid option!")
		}
		// small delay added
		time.Sleep(800 * time.Millisecond)
	}
}

// validateMainMenuInput validates input from user
func validateMainMenuInput(input int) bool {
	return input >= 1 && input <= 7
}

func createHireMenu() {

	hireData := getHireData()

	// create hire
	req := &hirepb.CreateHireRequest{
		Data: &hirepb.HireData{
			Name:     hireData.Name,
			Type:     hireData.Type,
			Duration: hireData.Duration,
			Role:     hireData.Role,
			Tags:     hireData.Tags,
		},
	}

	createHire(req)

}

// validateHireType makes the conversion from a string to hire type EMPLOYEE or CONTRACTOR
func validateHireType(str string) (hirepb.HireType, error) {
	var hire_type hirepb.HireType = 0

	if str == "employee" {
		hire_type = 1
	} else if str == "contractor" {
		hire_type = 2
	}

	if hire_type == 0 {
		return 0, fmt.Errorf("Employe type must be Employee or Contractor")
	}

	return hire_type, nil
}

// createHire handles the request to add a new hire to the DB
func createHire(req *hirepb.CreateHireRequest) {
	resp, err := ServiceClient.CreateHire(context.Background(), req)
	if err != nil {
		log.Println("Error creating hire data", err)
	}

	fmt.Println("\nHire added successully!!")
	fmt.Printf("Welcome aboard %s!!\n", resp.GetData().GetName())
}

// listHiresMenu handles both view and DB request for current hires
func listHiresMenu() {

	listStream, err := ServiceClient.ListHires(context.Background(), &hirepb.ListHireRequest{})

	if err != nil {
		log.Println("\nError listing blogs", err)
	} else {
		printStreamHeader()
		for {
			resp, recv_err := listStream.Recv()
			if recv_err != nil {
				if recv_err == io.EOF {
					break
				}
				log.Println("listHiresMenu RPC Stream receive error", recv_err)
			}

			if resp != nil {
				fmt.Printf("%-20s%-15v%-25s%-8d", resp.GetData().GetName(),
					resp.GetData().GetType(), resp.GetData().GetRole(),
					resp.GetData().GetDuration())
				for _, x := range resp.GetData().GetTags() {
					fmt.Printf("%s ", x)
				}
				fmt.Println()
			}
		}
	}
}

// findHiresMenu handles menu to look for hires.
// It can currently look for name, role, duration, and tags
func findHiresMenu() {
	var user_input int
	var find_pattern string
	var find_text string
	req := hirepb.FindHireRequest{
		FindPattern: find_pattern,
		FindText:    find_text,
	}

	for {
		fmt.Printf("%s", findMenu)
		fmt.Scanf("%d", &user_input)

		scanner := bufio.NewScanner(os.Stdin)

		if validateMainMenuInput(user_input) {
			switch user_input {
			case 1:
				req.FindPattern = "name"
				fmt.Printf("%s", findTxt)
				scanner.Scan()
				buffer := scanner.Text()
				req.FindText = strings.TrimSuffix(buffer, "\n")
			case 2:
				fmt.Println("Not yet supported")
				continue
			case 3:
				req.FindPattern = "role"
				fmt.Printf("%s", findTxt)
				scanner.Scan()
				buffer := scanner.Text()
				req.FindText = strings.TrimSuffix(buffer, "\n")
			case 4:
				req.FindPattern = "duration"
				fmt.Printf("%s", findTxt)
				fmt.Scanf("%d", &req.FindNumber)
			case 5:
				req.FindPattern = "tags"
				fmt.Printf("%s", findTxt)
				scanner.Scan()
				buffer := scanner.Text()
				req.FindText = strings.TrimSuffix(buffer, "\n")
			case 6:
				return
			}

		} else {
			fmt.Println("Invalid option!")
		}

		if req.FindPattern != "" {
			break
		}

		// small delay added
		time.Sleep(800 * time.Millisecond)
	}

	listStream, err := ServiceClient.FindHire(context.Background(), &req)

	if err != nil {
		log.Println("\nError listing blogs", err)
	} else {
		printStreamHeader()
		for {
			resp, recv_err := listStream.Recv()
			if recv_err != nil {
				if recv_err == io.EOF {
					fmt.Println("Search ended")
					break
				}
				log.Println("listHiresMenu RPC Stream receive error", recv_err)
			}

			if resp != nil {
				fmt.Printf("%-20s%-15v%-25s%-8d", resp.GetData().GetName(),
					resp.GetData().GetType(), resp.GetData().GetRole(),
					resp.GetData().GetDuration())
				for _, x := range resp.GetData().GetTags() {
					fmt.Printf("%s ", x)
				}
				fmt.Println()
			}
		}
	}
}

// deleteHire takes care of hire deletion from the DB
func deleteHire() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s", deleteTxt)
	scanner.Scan()
	del_name := strings.TrimSuffix(scanner.Text(), "\n")

	// check if hire exist and ask for confirmation
	find_req := hirepb.FindOneHireRequest{
		HireName: del_name,
	}
	find_res, err := ServiceClient.FindOneHire(context.Background(), &find_req)
	if err != nil {
		log.Println("deleteHire error:", err)
	}

	if find_res.GetFound() == true {
		fmt.Printf("Hire %s found.\n", del_name)
		fmt.Printf("%s", confirmDeleteTxt)
		scanner.Scan()
		buffer := strings.TrimSuffix(scanner.Text(), "\n")

		if buffer == "Y" || buffer == "y" {
			// make delete request to DB
			del_req := hirepb.DeleteHireRequest{
				HireName: del_name,
			}

			resp, err := ServiceClient.DeleteHire(context.Background(), &del_req)
			if err != nil {
				log.Println("\nCould not find hire", err)
			} else {
				fmt.Printf("\n%s deleted from database successully!!\n", resp.GetHireName())
			}
		}
	} else {
		fmt.Println("Hire not found")
	}
}

// updateHireData updates data for current hires
func updateHireData() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s", updateTxt)
	scanner.Scan()
	del_name := strings.TrimSuffix(scanner.Text(), "\n")

	// check if hire exist and ask for confirmation
	find_req := hirepb.FindOneHireRequest{
		HireName: del_name,
	}
	find_res, err := ServiceClient.FindOneHire(context.Background(), &find_req)
	if err != nil {
		log.Println("deleteHire error:", err)
	}

	if find_res.GetFound() == true {

		fmt.Printf("Hire %s found.\n", del_name)
		new_data := getHireData()

		new_data.ID, err = primitive.ObjectIDFromHex(find_res.GetData().GetId())
		if err != nil {
			log.Println("Cannot parse hire ID")
			return
		}

		req := hirepb.UpdateHireRequest{
			Data: model.DataToHirepb(new_data),
		}

		resp, err := ServiceClient.UpdateHire(context.Background(), &req)
		if err != nil {
			log.Println("updateHire error:", err)
		} else {
			fmt.Printf("Hire %s update succesfully\n", resp.GetData().GetName())
		}

	} else {
		fmt.Println("Hire not found")
	}
}

// getHireData obtains data from the user terminal
func getHireData() *model.HireDataItem {

	scanner := bufio.NewScanner(os.Stdin)
	input := &model.HireDataItem{}
	var buffer string

	for {
		// Obtains and validates hire type
		fmt.Printf("%s", createHireType)
		scanner.Scan()
		buffer = strings.TrimSuffix(scanner.Text(), "\n")
		buffer = strings.ToLower(buffer)
		hire_type, err := validateHireType(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		input.Type = hire_type

		// obtains name
		fmt.Printf("%s", createHireName)
		scanner.Scan()
		buffer = scanner.Text()
		input.Name = strings.TrimSuffix(buffer, "\n")

		// handles between CONTRACT and EMPLOYEE hire
		if input.Type == hirepb.HireType_EMPLOYEE {
			fmt.Printf("%s", createHireRole)
			scanner.Scan()
			buffer = scanner.Text()
			input.Role = strings.TrimSuffix(buffer, "\n")
			input.Duration = 999
		} else {
			fmt.Printf("%s", createHireDuration)
			fmt.Scanf("%d", &input.Duration)
			input.Role = "Contractor"
		}

		// handled tags
		fmt.Printf("%s", createHireTags)
		scanner.Scan()
		buffer = scanner.Text()
		input.Tags = strings.Split(buffer, ",")
		break
	}

	return input
}

// printStreamHeader is an utility function to print the header for hire data
func printStreamHeader() {
	fmt.Println("\nListing all hires in database")
	fmt.Printf("%-20s%7s%15s%23s%10s\n", "Name", "Type", "Role", "Duration", "Tags")
}

// addDummyEmployees currently adds dummy employees for demonstratin purposes
func addDummyEmployees() {
	var data = []hirepb.HireData{
		hirepb.HireData{Name: "John Smith", Type: hirepb.HireType_CONTRACTOR,
			Duration: 4, Role: "Developer", Tags: []string{"C++", "C#"}},
		hirepb.HireData{Name: "Peter Parker", Type: hirepb.HireType_EMPLOYEE,
			Duration: 99, Role: "Researcher", Tags: []string{"C++", "Golang"}},
		hirepb.HireData{Name: "Trott Irwin", Type: hirepb.HireType_CONTRACTOR,
			Duration: 3, Role: "Artist", Tags: []string{"Maya", "Blender"}},
		hirepb.HireData{Name: "George Brown", Type: hirepb.HireType_CONTRACTOR,
			Duration: 1, Role: "Painter", Tags: []string{"C++", "VBA", "Rust"}},
		hirepb.HireData{Name: "Xiang White", Type: hirepb.HireType_CONTRACTOR,
			Duration: 4, Role: "Developer", Tags: []string{"C++", "C#", "Embedded"}},
		hirepb.HireData{Name: "Blaise Pascal", Type: hirepb.HireType_EMPLOYEE,
			Duration: 99, Role: "Researcher", Tags: []string{"C++", "Golang", "Haskell"}},
		hirepb.HireData{Name: "Kiara Kovu", Type: hirepb.HireType_CONTRACTOR,
			Duration: 3, Role: "Developer", Tags: []string{"C", "Linux", "Erlang"}},
		hirepb.HireData{Name: "Nala Doc", Type: hirepb.HireType_CONTRACTOR,
			Duration: 4, Role: "Artist", Tags: []string{"Maya", "Linux"}},
		hirepb.HireData{Name: "Enrico Fermi", Type: hirepb.HireType_EMPLOYEE,
			Duration: 99, Role: "Researcher", Tags: []string{"C++", "Golang"}},
		hirepb.HireData{Name: "Pumba Brown", Type: hirepb.HireType_CONTRACTOR,
			Duration: 1, Role: "Researcher", Tags: []string{"C++", "VBA", "Linux"}},
		hirepb.HireData{Name: "Valdo Timon", Type: hirepb.HireType_CONTRACTOR,
			Duration: 3, Role: "Analyst", Tags: []string{"C++", "VBA"}},
	}

	// create hire
	for x, _ := range data {

		req := &hirepb.CreateHireRequest{
			Data: &hirepb.HireData{
				Name:     data[x].Name,
				Type:     data[x].Type,
				Duration: data[x].Duration,
				Role:     data[x].Role,
				Tags:     data[x].Tags,
			},
		}

		createHire(req)
	}

}
