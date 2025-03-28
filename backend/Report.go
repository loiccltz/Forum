package backend

import (
    "database/sql"
    "time"
)

type Report struct {
    ID            int
    ReporterID    int
    PostID        int
    CommentID     int
    Reason        string
    Status        string
    CreatedAt     time.Time
    ResolvedAt    *time.Time
    ResolvedByID  *int
}

const (
    StatusPending   = "pending"
    StatusApproved  = "approved"
    StatusRejected  = "rejected"
)

func CreatePostReport(db *sql.DB, reporterID, postID int, reason string) error {
    _, err := db.Exec(`
        INSERT INTO post_reports 
        (reporter_id, post_id, reason, status) 
        VALUES (?, ?, ?, ?)
    `, reporterID, postID, reason, StatusPending)
    return err
}

func CreateCommentReport(db *sql.DB, reporterID, commentID int, reason string) error {
    _, err := db.Exec(`
        INSERT INTO comment_reports 
        (reporter_id, comment_id, reason, status) 
        VALUES (?, ?, ?, ?)
    `, reporterID, commentID, reason, StatusPending)
    return err
}

func GetPendingReports(db *sql.DB) ([]Report, error) {
    rows, err := db.Query(`
        SELECT id, reporter_id, post_id, comment_id, reason, status, created_at 
        FROM post_reports 
        WHERE status = ?
    `, StatusPending)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reports []Report
    for rows.Next() {
        var report Report
        err := rows.Scan(
            &report.ID, 
            &report.ReporterID, 
            &report.PostID, 
            &report.CommentID, 
            &report.Reason, 
            &report.Status, 
            &report.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        reports = append(reports, report)
    }
    return reports, nil
}

func ResolveReport(db *sql.DB, reportID, resolverID int, approve bool) error {
    status := StatusRejected
    if approve {
        status = StatusApproved
    }

    _, err := db.Exec(`
        UPDATE post_reports 
        SET status = ?, 
            resolved_at = NOW(), 
            resolved_by_id = ? 
        WHERE id = ?
    `, status, resolverID, reportID)
    return err
}

func DeleteReportedContent(db *sql.DB, reportID int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    var postID, commentID *int
    err = tx.QueryRow(`
        SELECT post_id, comment_id 
        FROM post_reports 
        WHERE id = ?
    `, reportID).Scan(&postID, &commentID)
    if err != nil {
        tx.Rollback()
        return err
    }

    if postID != nil {
        _, err = tx.Exec("DELETE FROM post WHERE id = ?", *postID)
    } else if commentID != nil {
        _, err = tx.Exec("DELETE FROM comment WHERE id = ?", *commentID)
    }

    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`
        UPDATE post_reports 
        SET status = ?, 
            resolved_at = NOW() 
        WHERE id = ?
    `, StatusRejected, reportID)

    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}