import { useEffect, useState } from 'react';

import {
    Stack,
    Table,
    Container,
    TableBody,
    TableRow,
    TableCell,
    TableContainer,
    TablePagination,
    Backdrop,
    Typography,
    CircularProgress,
    Card,
} from '@mui/material';

import Scrollbar from 'src/components/scrollbar';

import SearchUserForm from 'src/components/search-user/searchUserForm';
import ErrorToast from 'src/components/toast/errorToast';
import SuccessToast from 'src/components/toast/successToast';

import TableNoData from '../table-no-data';

import RepoTableRow from '../search-table-row';
import RepoTableHead from '../search-table-head';
import RepoTableToolbar from '../search-table-toolbar';

import IssueTableRow from '../issues-table-row';
import IssueTableHead from '../issues-table-head';
import IssueTableToolbar from '../issues-table-toolbar';

import TableEmptyRows from '../table-empty-rows';
import { emptyRows, applyFilter, getComparator } from '../utils';

import {
    searchGithubUsers,
    searchGithubUsersRepo,
    saveRepos,
    searchRepoLinksForUser,
} from 'src/api/client';
import UserGrid from 'src/components/search-user/searchUserGrid';

// ----------------------------------------------------------------------

export default function SearchUserPage() {
    const [page, setPage] = useState(0);
    const [order, setOrder] = useState('asc');
    const [selected, setSelected] = useState([]);
    const [orderBy, setOrderBy] = useState('name');
    const [filterName, setFilterName] = useState('');
    const [rowsPerPage, setRowsPerPage] = useState(10);

    const [errorToast, setErrorToast] = useState({ open: false, message: '' });
    const [successToast, setSuccessToast] = useState({ open: false, message: '' });

    const [repos, setRepos] = useState([]);
    const [issues, setIssues] = useState([]);
    const [users, setUsers] = useState([]);
    const [username, setUsername] = useState('');
    const [selectedUser, setSelectedUser] = useState(null);

    // State mangagement for user flow
    const [isSearching, setIsSearching] = useState(false);
    const [isFindingRepos, setIsFindingRepos] = useState(false);
    const [isProcessingRepos, setIsProcessingRepos] = useState(false);

    const [hasSearched, setHasSearched] = useState(false);
    const [hasPickedUser, setHasPickedUser] = useState(false);
    const [hasProcessedRepos, setHasProcessedRepos] = useState(false);

    useEffect(() => { }, []);

    const handleSort = (event, id) => {
        const isAsc = orderBy === id && order === 'asc';
        if (id !== '') {
            setOrder(isAsc ? 'desc' : 'asc');
            setOrderBy(id);
        }
    };

    const handleSelectAllClick = (event) => {
        if (event.target.checked) {
            const newSelecteds = repos.map((n) => n);
            setSelected(newSelecteds);
            return;
        }
        setSelected([]);
    };

    const handleClick = (event, row, index) => {
        const selectedIndex = selected.findIndex((item) => item === row);
        let newSelected = [];

        if (selectedIndex === -1) {
            // If the row is not in the selected array, add it
            newSelected = [...selected, row];
        } else if (selectedIndex === 0) {
            // If it's the first item, remove it
            newSelected = selected.slice(1);
        } else if (selectedIndex === selected.length - 1) {
            // If it's the last item, remove it
            newSelected = selected.slice(0, -1);
        } else if (selectedIndex > 0) {
            // If it's in the middle, remove it
            newSelected = [
                ...selected.slice(0, selectedIndex),
                ...selected.slice(selectedIndex + 1),
            ];
        }

        setSelected(newSelected);
    };

    const handleChangePage = (event, newPage) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event) => {
        setPage(0);
        setRowsPerPage(parseInt(event.target.value, 10));
    };

    const handleFilterByName = (event) => {
        setPage(0);
        setFilterName(event.target.value);
    };

    const dataFiltered = applyFilter({
        inputData: repos,
        comparator: getComparator(order, orderBy),
        filterName,
    });

    const issuesFiltered = applyFilter({
        inputData: issues,
        comparator: getComparator(order, orderBy),
        filterName,
    });

    const handleCloseErrorToast = () => {
        setErrorToast({ open: false, message: '' });
    };

    const handleCloseSuccessToast = () => {
        setSuccessToast({ open: false, message: '' });
    };

    const handleSearch = async () => {
        if (isSearching || isFindingRepos) {
            return;
        }
        setIsSearching(true);
        setHasSearched(false);
        setHasPickedUser(false);
        setHasProcessedRepos(false);
        console.log(`Starting github search for ${username}`);

        try {
            const res = await searchGithubUsers(username);
            if (res === null) {
                setUsers([]);
            } else {
                setUsers(res);
            }

            setIsSearching(false);
            setHasSearched(true);
        } catch (e) {
            setErrorToast({ open: true, message: e.response.data.message });
        }
        setIsSearching(false);
    };

    const handleUserSelect = async (login) => {
        setSelectedUser(login);
        console.log(`Selected user: ${login}..getting repos`);
        setIsFindingRepos(true);

        //  Find Repos for user
        var repositories;
        try {
            const res = await searchGithubUsersRepo(login);
            if (res === null) {
                setRepos([]);
            } else {
                setRepos(res);
                repositories = res;
            }
            // setIsFindingRepos(false);
            setHasSearched(false);
            setHasPickedUser(true);
        } catch (e) {
            setErrorToast({ open: true, message: e.response.data.message });
            setIsFindingRepos(false);
            setHasSearched(false);
            setHasPickedUser(true);
        }

        // Save repos to DB
        try {
            await saveRepos(repositories);
            console.log(repositories.length, 'repos saved to db');
        } catch (e) {
            console.log(e);
            setErrorToast({ open: true, message: e.response.data.message });
        }

        // Pull states of repos
        try {
            const res = await searchRepoLinksForUser(login);
            if (res === null) {
                setIssues([]);
            } else {
                setIssues(res);
            }

            console.log('Search completed!');
            setHasProcessedRepos(true);
            setHasSearched(false);
            setHasPickedUser(false);
        } catch (e) {
            setErrorToast({ open: true, message: e.response.data.message });
            setHasProcessedRepos(true);
            setHasSearched(false);
            setHasPickedUser(false);
        }
        setSuccessToast({
            open: true,
            message: `Found ${repositories.length} repo. Please wait a few seconds to see results.`,
        });
    };

    // The save button
    const handleProcessRepos = async () => {
        // if (isProcessingRepos | isFindingRepos) {
        //     return;
        // }
        setIsProcessingRepos(true);
        console.log('Saving & Processing repos');
        console.log('Saving repos to db');
        console.log('Repos to save:', repos);

        // Step 1. We save the repos to DB.
        try {
            await saveRepos(repos);
            console.log(repos.length, 'repos saved to db');
            // setSuccessToast({
            //     open: true,
            //     message: `${selected.length} repo(s) queued for processing. Please wait a few seconds.`,
            // });
        } catch (e) {
            console.log(e);
            setErrorToast({ open: true, message: e.response.data.message });
        }
        setSelected([]);

        // Step 2. We pull all repo links for user (eventually they will be populated by the parser job)
        try {
            const res = await searchRepoLinksForUser(selectedUser);
            if (res === null) {
                setIssues([]);
            } else {
                setIssues(res);
            }

            console.log('Search completed!');
            setHasProcessedRepos(true);
            setHasSearched(false);
            setHasPickedUser(false);
        } catch (e) {
            setErrorToast({ open: true, message: e.response.data.message });
            setHasProcessedRepos(true);
            setHasSearched(false);
            setHasPickedUser(false);
        }
        setSuccessToast({
            open: true,
            message: `${selected.length} repo(s) queued for processing. Please wait a few seconds.`,
        });
    };

    const notFound = !dataFiltered.length && !!filterName;

    return (
        <Container>
            <Stack direction="column" alignItems="left" justifyContent="space-between" mb={5}>
                <Stack mb={1}>
                    <Typography variant="h4">Search for new Repositories</Typography>
                </Stack>
                <Typography variant="p">
                    Here you can search GitHub for new repositories to add to the GooGL-Bye queue.
                </Typography>
            </Stack>

            <SearchUserForm
                setUsername={setUsername}
                handleSearch={handleSearch}
                username={username}
                isSearching={isSearching}
            />

            {hasSearched && users.length > 0 ? (
                <UserGrid
                    isFindingRepos={isFindingRepos}
                    users={users}
                    onUserSelect={handleUserSelect}
                />
            ) : hasSearched ? (
                <Card sx={{ mt: 4, p: 3, textAlign: 'center' }}>
                    <Typography variant="h6" color="text.secondary">
                        No users found
                    </Typography>
                    <Typography variant="body2" sx={{ mt: 1 }}>
                        We couldn't find any GitHub users matching your search. Please try a
                        different username.
                    </Typography>
                </Card>
            ) : null}

            {hasProcessedRepos && (
                <Card>
                    <IssueTableToolbar
                        numSelected={selected.length}
                        filterName={filterName}
                        onFilterName={handleFilterByName}
                    />

                    <Scrollbar>
                        <TableContainer sx={{ overflow: 'unset' }}>
                            <Table sx={{ minWidth: 800 }}>
                                <IssueTableHead
                                    order={order}
                                    orderBy={orderBy}
                                    rowCount={issues.length}
                                    numSelected={selected.length}
                                    onRequestSort={handleSort}
                                    headLabel={[
                                        { id: 'name', label: 'Repository Name' },
                                        { id: 'author', label: 'Author' },
                                        { id: 'language', label: 'Language' },
                                        { id: 'stars', label: 'Stars', align: 'center' },
                                        { id: 'forks', label: 'Forks', align: 'center' },
                                        { id: 'state', label: 'Status' },
                                        { id: 'links', label: 'Links' },
                                        { id: '' },
                                    ]}
                                />
                                <TableBody>
                                    {issues.length > 0 ? (
                                        <>
                                            {issuesFiltered
                                                .slice(
                                                    page * rowsPerPage,
                                                    page * rowsPerPage + rowsPerPage
                                                )
                                                .map((row) => (
                                                    <IssueTableRow
                                                        key={row.id}
                                                        name={row.name}
                                                        author={row.author}
                                                        language={row.language}
                                                        stars={row.stars}
                                                        forks={row.forks}
                                                        ghUrl={row.ghUrl}
                                                        state={row.state}
                                                        avatarUrl={row.avatarUrl}
                                                        lastCommit={row.lastCommit}
                                                        issues={row.links}
                                                        errorMsg={row.errorMsg}
                                                        selected={selected.indexOf(row.name) !== -1}
                                                    />
                                                ))}
                                            <TableEmptyRows
                                                height={77}
                                                emptyRows={emptyRows(
                                                    page,
                                                    rowsPerPage,
                                                    issues.length
                                                )}
                                            />
                                            {notFound && <TableNoData query={filterName} />}
                                        </>
                                    ) : (
                                        <TableRow>
                                            <TableCell colSpan={7} align="center">
                                                <Typography variant="subtitle1">
                                                    Nice! No goo.gl links found in repos for user{' '}
                                                    {selectedUser}.
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Scrollbar>
                    <TablePagination
                        page={page}
                        component="div"
                        count={repos.length}
                        rowsPerPage={rowsPerPage}
                        onPageChange={handleChangePage}
                        rowsPerPageOptions={[5, 10, 25]}
                        onRowsPerPageChange={handleChangeRowsPerPage}
                    />
                </Card>
            )}

            <ErrorToast
                open={errorToast.open}
                message={errorToast.message}
                onClose={handleCloseErrorToast}
            />
            <SuccessToast
                open={successToast.open}
                message={successToast.message}
                onClose={handleCloseSuccessToast}
            />
        </Container>
    );
}
