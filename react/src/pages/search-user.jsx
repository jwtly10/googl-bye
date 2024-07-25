import { Helmet } from 'react-helmet-async';

import { SearchUserView } from 'src/sections/search-user/view';

// ----------------------------------------------------------------------

export default function SearchUserPage() {
    return (
        <>
            <Helmet>
                <title> Search-User | GooGl-Bye </title>
            </Helmet>

            <SearchUserView />
        </>
    );
}
